//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package service

import (
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/api/client"
)

const logLevel = 2

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("api:handler:service:list list services in %s", nid)

	var (
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:list get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:service:list get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	items, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:list get service list in namespace `%s` err: %s", ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	dl, err := dm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:list get pod list by service id `%s` err: %s", ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().NewList(items, dl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:list convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("api:handler:service:list write response err: %s", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("api:handler:service:info get service `%s` in namespace `%s`", sid, nid)

	var (
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:info get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:service:info get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:info get service by name `%s` in namespace `%s` err: %s", sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("api:handler:service:info service `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dl, err := dm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:info get pod list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	pods, err := pdm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:info get pod list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().New(srv, dl, pods).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:info convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("handler:service:info write response err: %s", err.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("handler:service:create create service in namespace `%s`", nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Service().CreateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("handler:service:create validation incoming data err: %s", e.Err())
		e.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:create get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:service:create get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	if opts.Name == nil {
		data, err := converter.DockerNamespaceParse(*opts.Image)
		if err != nil {
			errors.New("service").BadParameter("image").Http(w)
			return
		}
		opts.Name = &data.Repo
	}

	srv, err := sm.Get(ns.Meta.Name, *opts.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:create get service by name `%s` in namespace `%s` err: %s", opts.Name, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv != nil {
		log.V(logLevel).Warnf("api:handler:service:create service name `%s` in namespace `%s` not unique", opts.Name, ns.Meta.Name)
		errors.New("service").NotUnique("name").Http(w)
		return
	}

	srv, err = sm.Create(ns, opts)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:create create service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().New(srv, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:create convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("api:handler:service:create write response err: %s", err.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("api:handler:service:update update service `%s` in namespace `%s`", sid, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Service().UpdateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("api:handler:service:update validation incoming data err: %s", e.Err())
		e.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:update get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:service:update get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service: get service by name` err: %s", sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("api:handler:service:update service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	srv, err := sm.Update(svc, opts)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:update update service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().New(srv, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:update convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("handler:service:update write response err: %s", err.Error())
		return
	}
}

func ServiceLogsH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]
	did := r.URL.Query().Get("deployment")
	pid := r.URL.Query().Get("pod")
	cid := r.URL.Query().Get("container")
	notify := w.(http.CloseNotifier).CloseNotify()
	doneChan := make(chan bool, 1)

	go func() {
		<-notify
		log.Debug("api:handler:service:logs HTTP connection just closed.")
		doneChan <- true
	}()

	log.V(logLevel).Debugf("api:handler:service:logs get logs service `%s` in namespace `%s`", sid, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		pm  = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:service:logs get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs get service by name` err: %s", sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("api:handler:service:logs service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	deployment, err := dm.Get(ns.Meta.Name, svc.Meta.Name, did)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs get deployment by name` err: %s", did, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if deployment == nil {
		log.V(logLevel).Warnf("api:handler:service:logs deployment `%s` not found", pid)
		errors.New("service").NotFound().Http(w)
		return
	}

	pod, err := pm.Get(ns.Meta.Name, did, pid, cid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs get pod by name` err: %s", sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("api:handler:service:logs pod `%s` not found", pid)
		errors.New("service").NotFound().Http(w)
		return
	}

	node, err := nm.Get(pod.Meta.Node)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs get node by name err: %s", sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if node == nil {
		log.V(logLevel).Warnf("api:handler:service:logs node %s not found", sid, pod.Meta.Node)
		errors.New("service").NotFound().Http(w)
		return
	}

	httpcli, err := client.NewHTTP("http://"+node.Info.InternalIP, &client.Config{BearerToken: ""})
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs create http client err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	res, err := httpcli.V1().Cluster().Node().Logs(r.Context(), pid, cid, nil)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:logs get pod logs err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	var buffer []byte

	go func() {
		for {
			select {
			case <-doneChan:
				res.Close()
				return
			default:
				n, err := res.Read(buffer)
				if err != nil {
					log.V(logLevel).Errorf("api:handler:service:logs read bytes from stream err: %s", err)
					res.Close()
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = w.Write(p)
					if err != nil {
						log.V(logLevel).Errorf("api:handler:service:logs write bytes from stream err: %s", err)
						return n, err
					}
					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])
				if err != nil {
					log.V(logLevel).Errorf("api:handler:service:logs written to stream err: %s", err)
					return
				}

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}
			}
		}
	}()

}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("handler:service:remove remove service `%s` from app `%s`", sid, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("handler:service:remove get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:service:remove get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:service:remove get service by name `%s` in namespace `%s` err: %s", sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("api:handler:service:remove service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	if _, err := sm.Destroy(svc); err != nil {
		log.V(logLevel).Errorf("api:handler:service:remove remove service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("api:handler:service:remove write response err: %s", err.Error())
		return
	}
}
