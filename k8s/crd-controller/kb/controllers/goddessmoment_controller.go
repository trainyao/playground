/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kbv1 "github.com/trainyao/playground/k8s/crd-controller/kb/api/v1"
)

var (
	foodFocus = map[string]string{
		"小明": "珍珠奶茶",
		"小王": "麻辣烫",
		"小杰": "螺蛳粉",
	}
	buyFoodTimeCostMap = map[string]time.Duration{
		"珍珠奶茶": time.Second,
		"麻辣烫":  5 * time.Second,
		"螺蛳粉":  10 * time.Second,
	}
)

// GoddessMomentReconciler reconciles a GoddessMoment object
type GoddessMomentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Name   string
}

//+kubebuilder:rbac:groups=kb.crd.playground.trainyao.io,resources=goddessmoments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kb.crd.playground.trainyao.io,resources=goddessmoments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kb.crd.playground.trainyao.io,resources=goddessmoments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GoddessMoment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *GoddessMomentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	res := ctrl.Result{}

	var goddessMoment kbv1.GoddessMoment
	err := r.Get(ctx, req.NamespacedName, &goddessMoment)
	if err != nil {
		// The GoddessMoment resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			return res, fmt.Errorf("goddessMoment '%s' in work queue no longer exists", req.Name)
		}

		res.Requeue = true
		return res, err
	}

	focusFood, ok := foodFocus[r.Name]
	if !ok {
		return res, fmt.Errorf("%s 朋友圈找不到 %s 的关注食物", goddessMoment.Name, r.Name)
	}
	logger.Info("关注食物", "name", r.Name, "food", focusFood)

	foundFocusFood := false
	for _, food := range goddessMoment.Spec.FoodDemand {
		if food.Name == focusFood {
			foundFocusFood = true
			break
		}
	}

	if !foundFocusFood {
		logger.Info("关注食物没有在朋友圈发现", "name", r.Name, "food", focusFood)
		return res, nil
	}

	var foodStatus *kbv1.FoodDemandStatus
	var index = -1
	for i, foodDemandStatus := range goddessMoment.Status.FoodDemand {
		if foodDemandStatus.Name == focusFood {
			foodStatus = &goddessMoment.Status.FoodDemand[i]
			index = i
			break
		}
	}

	if foodStatus == nil {
		return res, fmt.Errorf("关注食物 %s 在女神朋友圈出现, 但是在status里面没有找到", focusFood)
	}
	logger.Info("朋友圈关注的食物状态是", "name", r.Name, "food", focusFood, "status", foodStatus.Status)

	if foodStatus.Status == kbv1.FoodStatusPending {
		// 处理女神朋友圈里关注食物还未被人认领的情况
		err = r.handlePending(ctx, &goddessMoment, index)
	} else {
		// 处理关注食物已经有人认领的情况
		// 认领的人不是我
		if foodStatus.ClaimBy != r.Name {
			logger.Info("朋友圈食物已经被人认领了, 我来的太迟了", "food", focusFood, "name", foodStatus.ClaimBy)
			return res, nil
		}
		if foodStatus.Status == kbv1.FoodStatusPendingArrival {
			// 去购买食物
			err = r.handlePendingArrival(ctx, &goddessMoment, index)
		}

		// 状态是已到达, 本条朋友圈任务已完成
		logger.Info("朋友圈食物我已经认领并送到了, 可以不用关注了", "food", focusFood)
	}
	if err != nil {
		err = fmt.Errorf("处理朋友圈更新 %s 失败, err: %s", goddessMoment.Name, err.Error())
		res.Requeue = true
		return res, err
	}

	return res, nil
}

func (r *GoddessMomentReconciler) handlePending(ctx context.Context, gm *kbv1.GoddessMoment, index int) error {
	logger := log.FromContext(ctx)

	focusFood := gm.Status.FoodDemand[index].Name
	logger.Info("开始告诉女神我来送", "food", focusFood)
	if r.Name == "小杰" {
		data := `[{"op":"replace","path":"/status/foodDemand","value":%s}]`
		foodDemandStatusCopy := gm.Status.DeepCopy()
		foodDemandStatusCopy.FoodDemand[index].Status = kbv1.FoodStatusPendingArrival
		foodDemandStatusCopy.FoodDemand[index].ClaimTime = metav1.Now()
		foodDemandStatusCopy.FoodDemand[index].ClaimBy = r.Name
		valueJson, err := json.Marshal(&foodDemandStatusCopy.FoodDemand)
		if err != nil {
			return fmt.Errorf("error when marshal foodDemand status, err: %s", err.Error())
		}

		dataJson := fmt.Sprintf(data, valueJson)
		err = r.Patch(ctx, gm, client.RawPatch(types.JSONPatchType, []byte(dataJson)))
		if err != nil {
			return fmt.Errorf("patch更新朋友圈告诉女神 %s 我来送 %s 失败, err: %s, datajson: %s",
				gm.Name, focusFood, err.Error(), dataJson)
		}
	} else {
		momentCopy := gm.DeepCopy()
		momentCopy.Status.FoodDemand[index].Status = kbv1.FoodStatusPendingArrival
		momentCopy.Status.FoodDemand[index].ClaimTime = metav1.Now()
		momentCopy.Status.FoodDemand[index].ClaimBy = r.Name

		err := r.Update(ctx, momentCopy)
		if err != nil {
			return fmt.Errorf("更新朋友圈 %s 告诉女神 %s 我来送 失败, err: %s",
				gm.Name, focusFood, err.Error())
		}
	}
	logger.Info("更新朋友圈告诉女神送来成功", "food", focusFood)
	return nil
}

func (r *GoddessMomentReconciler) handlePendingArrival(ctx context.Context, gm *kbv1.GoddessMoment, index int) error {
	logger := log.FromContext(ctx)

	focusFood := gm.Status.FoodDemand[index].Name
	buyFoodTimeCost := buyFoodTimeCostMap[focusFood]
	logger.Info("为朋友圈购买需要秒, 正在前往购买", "food", focusFood, "buyFoodTime", buyFoodTimeCost/time.Second)
	time.Sleep(buyFoodTimeCost)
	logger.Info("为朋友圈购买完成, 更新女神朋友圈", "food", focusFood)

	momentCopy := gm.DeepCopy()
	momentCopy.Status.FoodDemand[index].Status = kbv1.FoodStatusArrived
	momentCopy.Status.FoodDemand[index].ArrivalTime = metav1.Now()

	err := r.Update(ctx, momentCopy)
	if err != nil {
		return fmt.Errorf("更新朋友圈 %s 失败, err: %s", gm.Name, err.Error())
	}
	logger.Info("更新朋友圈成功")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GoddessMomentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kbv1.GoddessMoment{}).
		Complete(r)
}
