package converter

import (
	"errors"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/resource"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/1.5/pkg/types"
	"k8s.io/client-go/1.5/pkg/util/intstr"
)

var (
	error_incoming_data  = errors.New("data incoming cannot be nil")
	error_outcoming_data = errors.New("data outcoming cannot be nil")
)

func Convert_v1_ObjectMeta_to_api_ObjectMeta(in *v1.ObjectMeta, out *api.ObjectMeta) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.Name = in.Name
	out.GenerateName = in.GenerateName
	out.Namespace = in.Namespace
	out.SelfLink = in.SelfLink
	out.UID = types.UID(in.UID)
	out.ResourceVersion = in.ResourceVersion
	out.Generation = in.Generation
	out.CreationTimestamp = in.CreationTimestamp
	out.DeletionTimestamp = in.DeletionTimestamp
	out.DeletionGracePeriodSeconds = in.DeletionGracePeriodSeconds
	out.Labels = in.Labels
	out.Annotations = in.Annotations
	out.Finalizers = in.Finalizers
	out.ClusterName = in.ClusterName

	if in.OwnerReferences != nil {
		in, out := &in.OwnerReferences, &out.OwnerReferences
		*out = make([]api.OwnerReference, len(*in))

		for i := range *in {
			if err := Convert_v1_OwnerReference_To_api_OwnerReference(&(*in)[i], &(*out)[i]); err != nil {
				return err
			}
		}
	} else {
		out.OwnerReferences = nil
	}

	return nil
}

func Convert_v1_PodTemplateSpec_to_api_PodTemplateSpec(in *v1.PodTemplateSpec, out *api.PodTemplateSpec) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.Name = in.Name
	out.GenerateName = in.GenerateName
	out.Namespace = in.Namespace
	out.SelfLink = in.SelfLink
	out.UID = in.UID
	out.ResourceVersion = in.ResourceVersion
	out.Generation = in.Generation
	out.CreationTimestamp = in.CreationTimestamp
	out.DeletionTimestamp = in.DeletionTimestamp
	out.DeletionGracePeriodSeconds = in.DeletionGracePeriodSeconds
	out.Labels = in.Labels
	out.Annotations = in.Annotations
	out.Finalizers = in.Finalizers
	out.ClusterName = in.ClusterName

	for _, val := range in.OwnerReferences {
		out.OwnerReferences = append(out.OwnerReferences, api.OwnerReference{
			APIVersion: val.APIVersion,
			Kind:       val.Kind,
			Name:       val.Name,
			UID:        val.UID,
			Controller: val.Controller,
		})
	}

	err := Convert_v1_ObjectMeta_to_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta)
	if err != nil {
		return err
	}

	err = Convert_v1_PodSpec_to_api_PodSpec(&in.Spec, &out.Spec)
	if err != nil {
		return err
	}

	return nil
}

func Convert_v1_PodSpec_to_api_PodSpec(in *v1.PodSpec, out *api.PodSpec) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	for _, val := range in.InitContainers {
		var container = new(api.Container)

		err := Convert_v1_Container_to_api_Container(&val, container)
		if err != nil {
			return err
		}

		out.InitContainers = append(out.InitContainers, *container)
	}

	for _, val := range in.Containers {

		var container = new(api.Container)

		err := Convert_v1_Container_to_api_Container(&val, container)
		if err != nil {
			return err
		}

		out.Containers = append(out.Containers, *container)
	}

	out.RestartPolicy = api.RestartPolicy(in.RestartPolicy)
	out.TerminationGracePeriodSeconds = in.TerminationGracePeriodSeconds
	out.ActiveDeadlineSeconds = in.ActiveDeadlineSeconds
	out.DNSPolicy = api.DNSPolicy(in.DNSPolicy)
	out.NodeSelector = in.NodeSelector
	out.ServiceAccountName = in.ServiceAccountName
	out.NodeName = in.NodeName

	if out.SecurityContext == nil {
		out.SecurityContext = new(api.PodSecurityContext)
	}

	out.SecurityContext.HostNetwork = in.HostNetwork
	out.SecurityContext.HostPID = in.HostPID
	out.SecurityContext.HostIPC = in.HostIPC

	if in.SecurityContext.SELinuxOptions != nil {
		if out.SecurityContext.SELinuxOptions == nil {
			out.SecurityContext.SELinuxOptions = new(api.SELinuxOptions)
		}
		out.SecurityContext.SELinuxOptions.Level = in.SecurityContext.SELinuxOptions.Level
		out.SecurityContext.SELinuxOptions.Role = in.SecurityContext.SELinuxOptions.Role
		out.SecurityContext.SELinuxOptions.Type = in.SecurityContext.SELinuxOptions.Type
		out.SecurityContext.SELinuxOptions.User = in.SecurityContext.SELinuxOptions.User
	}

	out.SecurityContext.RunAsUser = in.SecurityContext.RunAsUser
	out.SecurityContext.RunAsNonRoot = in.SecurityContext.RunAsNonRoot
	out.SecurityContext.SupplementalGroups = in.SecurityContext.SupplementalGroups
	out.SecurityContext.FSGroup = in.SecurityContext.FSGroup

	for _, val := range in.ImagePullSecrets {
		item := api.LocalObjectReference{Name: val.Name}
		out.ImagePullSecrets = append(out.ImagePullSecrets, item)
	}

	out.Hostname = in.Hostname
	out.Subdomain = in.Subdomain

	for _, val := range in.Volumes {
		var volume = api.Volume{}

		volume.Name = val.Name
		volume.VolumeSource = api.VolumeSource{}
		if val.VolumeSource.HostPath != nil {
			volume.VolumeSource.HostPath = new(api.HostPathVolumeSource)
			volume.VolumeSource.HostPath.Path = val.VolumeSource.HostPath.Path
		}

		if val.VolumeSource.EmptyDir != nil {
			volume.VolumeSource.EmptyDir = new(api.EmptyDirVolumeSource)
			volume.VolumeSource.EmptyDir.Medium = api.StorageMedium(val.VolumeSource.EmptyDir.Medium)
		}

		if val.VolumeSource.GCEPersistentDisk != nil {
			volume.VolumeSource.GCEPersistentDisk = new(api.GCEPersistentDiskVolumeSource)
			volume.VolumeSource.GCEPersistentDisk.FSType = val.VolumeSource.GCEPersistentDisk.FSType
			volume.VolumeSource.GCEPersistentDisk.Partition = val.VolumeSource.GCEPersistentDisk.Partition
			volume.VolumeSource.GCEPersistentDisk.PDName = val.VolumeSource.GCEPersistentDisk.PDName
			volume.VolumeSource.GCEPersistentDisk.ReadOnly = val.VolumeSource.GCEPersistentDisk.ReadOnly
		}

		if val.VolumeSource.AWSElasticBlockStore != nil {
			volume.VolumeSource.AWSElasticBlockStore = new(api.AWSElasticBlockStoreVolumeSource)
			volume.VolumeSource.AWSElasticBlockStore.ReadOnly = val.VolumeSource.AWSElasticBlockStore.ReadOnly
			volume.VolumeSource.AWSElasticBlockStore.Partition = val.VolumeSource.AWSElasticBlockStore.Partition
			volume.VolumeSource.AWSElasticBlockStore.FSType = val.VolumeSource.AWSElasticBlockStore.FSType
			volume.VolumeSource.AWSElasticBlockStore.VolumeID = val.VolumeSource.AWSElasticBlockStore.VolumeID
		}

		if val.VolumeSource.GitRepo != nil {
			volume.VolumeSource.GitRepo = new(api.GitRepoVolumeSource)
			volume.VolumeSource.GitRepo.Directory = val.VolumeSource.GitRepo.Directory
			volume.VolumeSource.GitRepo.Repository = val.VolumeSource.GitRepo.Repository
			volume.VolumeSource.GitRepo.Revision = val.VolumeSource.GitRepo.Repository
		}

		if val.VolumeSource.Secret != nil {
			volume.VolumeSource.Secret = new(api.SecretVolumeSource)
			volume.VolumeSource.Secret.DefaultMode = val.VolumeSource.Secret.DefaultMode
			volume.VolumeSource.Secret.SecretName = val.VolumeSource.Secret.SecretName

			for _, val := range val.VolumeSource.Secret.Items {
				volume.VolumeSource.Secret.Items = append(volume.VolumeSource.Secret.Items, api.KeyToPath{
					Key:  val.Key,
					Path: val.Path,
					Mode: val.Mode,
				})
			}
		}

		if val.VolumeSource.NFS != nil {
			volume.VolumeSource.NFS = new(api.NFSVolumeSource)
			volume.VolumeSource.NFS.Path = val.VolumeSource.NFS.Path
			volume.VolumeSource.NFS.ReadOnly = val.VolumeSource.NFS.ReadOnly
			volume.VolumeSource.NFS.Server = val.VolumeSource.NFS.Server
		}

		if val.VolumeSource.ISCSI != nil {
			volume.VolumeSource.ISCSI = new(api.ISCSIVolumeSource)
			volume.VolumeSource.ISCSI.ReadOnly = val.VolumeSource.ISCSI.ReadOnly
			volume.VolumeSource.ISCSI.FSType = val.VolumeSource.ISCSI.FSType
			volume.VolumeSource.ISCSI.IQN = val.VolumeSource.ISCSI.IQN
			volume.VolumeSource.ISCSI.ISCSIInterface = val.VolumeSource.ISCSI.ISCSIInterface
			volume.VolumeSource.ISCSI.Lun = val.VolumeSource.ISCSI.Lun
			volume.VolumeSource.ISCSI.TargetPortal = val.VolumeSource.ISCSI.TargetPortal
		}

		if val.VolumeSource.Glusterfs != nil {
			volume.VolumeSource.Glusterfs = new(api.GlusterfsVolumeSource)
			volume.VolumeSource.Glusterfs.ReadOnly = val.VolumeSource.Glusterfs.ReadOnly
			volume.VolumeSource.Glusterfs.EndpointsName = val.VolumeSource.Glusterfs.EndpointsName
			volume.VolumeSource.Glusterfs.Path = val.VolumeSource.Glusterfs.Path
		}

		if val.VolumeSource.PersistentVolumeClaim != nil {
			volume.VolumeSource.PersistentVolumeClaim = new(api.PersistentVolumeClaimVolumeSource)
			volume.VolumeSource.PersistentVolumeClaim.ReadOnly = val.VolumeSource.PersistentVolumeClaim.ReadOnly
			volume.VolumeSource.PersistentVolumeClaim.ClaimName = val.VolumeSource.PersistentVolumeClaim.ClaimName
		}

		if val.VolumeSource.RBD != nil {
			volume.VolumeSource.RBD = new(api.RBDVolumeSource)
			volume.VolumeSource.RBD.ReadOnly = val.VolumeSource.RBD.ReadOnly
			volume.VolumeSource.RBD.CephMonitors = val.VolumeSource.RBD.CephMonitors
			volume.VolumeSource.RBD.FSType = val.VolumeSource.RBD.FSType
			volume.VolumeSource.RBD.Keyring = val.VolumeSource.RBD.Keyring
			volume.VolumeSource.RBD.RadosUser = val.VolumeSource.RBD.RadosUser
			volume.VolumeSource.RBD.RBDImage = val.VolumeSource.RBD.RBDImage
			volume.VolumeSource.RBD.RBDPool = val.VolumeSource.RBD.RBDPool

			if val.VolumeSource.RBD.SecretRef != nil {
				volume.VolumeSource.RBD.SecretRef = new(api.LocalObjectReference)
				volume.VolumeSource.RBD.SecretRef.Name = val.VolumeSource.RBD.SecretRef.Name
			}
		}

		if val.VolumeSource.Quobyte != nil {
			volume.VolumeSource.Quobyte = new(api.QuobyteVolumeSource)
			volume.VolumeSource.Quobyte.Group = val.VolumeSource.Quobyte.Group
			volume.VolumeSource.Quobyte.ReadOnly = val.VolumeSource.Quobyte.ReadOnly
			volume.VolumeSource.Quobyte.Registry = val.VolumeSource.Quobyte.Registry
			volume.VolumeSource.Quobyte.User = val.VolumeSource.Quobyte.User
			volume.VolumeSource.Quobyte.Volume = val.VolumeSource.Quobyte.Volume
		}

		if val.VolumeSource.FlexVolume != nil {
			volume.VolumeSource.FlexVolume = new(api.FlexVolumeSource)
			volume.VolumeSource.FlexVolume.ReadOnly = val.VolumeSource.FlexVolume.ReadOnly
			volume.VolumeSource.FlexVolume.Driver = val.VolumeSource.FlexVolume.Driver
			volume.VolumeSource.FlexVolume.FSType = val.VolumeSource.FlexVolume.FSType
			volume.VolumeSource.FlexVolume.Options = val.VolumeSource.FlexVolume.Options

			if val.VolumeSource.FlexVolume.SecretRef != nil {
				volume.VolumeSource.FlexVolume.SecretRef = new(api.LocalObjectReference)
				volume.VolumeSource.FlexVolume.SecretRef.Name = val.VolumeSource.FlexVolume.SecretRef.Name
			}
		}

		if val.VolumeSource.Cinder != nil {
			volume.VolumeSource.Cinder = new(api.CinderVolumeSource)
			volume.VolumeSource.Cinder.ReadOnly = val.VolumeSource.Cinder.ReadOnly
			volume.VolumeSource.Cinder.FSType = val.VolumeSource.Cinder.FSType
			volume.VolumeSource.Cinder.VolumeID = val.VolumeSource.Cinder.VolumeID
		}

		if val.VolumeSource.CephFS != nil {
			volume.VolumeSource.CephFS = new(api.CephFSVolumeSource)
			volume.VolumeSource.CephFS.ReadOnly = val.VolumeSource.CephFS.ReadOnly
			volume.VolumeSource.CephFS.Monitors = val.VolumeSource.CephFS.Monitors
			volume.VolumeSource.CephFS.Path = val.VolumeSource.CephFS.Path
			volume.VolumeSource.CephFS.SecretFile = val.VolumeSource.CephFS.SecretFile
			volume.VolumeSource.CephFS.User = val.VolumeSource.CephFS.User

			if val.VolumeSource.CephFS.SecretRef != nil {
				volume.VolumeSource.CephFS.SecretRef = new(api.LocalObjectReference)
				volume.VolumeSource.CephFS.SecretRef.Name = val.VolumeSource.CephFS.SecretRef.Name
			}
		}

		if val.VolumeSource.Flocker != nil {
			volume.VolumeSource.Flocker = new(api.FlockerVolumeSource)
			volume.VolumeSource.Flocker.DatasetName = val.VolumeSource.Flocker.DatasetName
		}

		if val.VolumeSource.DownwardAPI != nil {
			volume.VolumeSource.DownwardAPI = new(api.DownwardAPIVolumeSource)
			volume.VolumeSource.DownwardAPI.DefaultMode = val.VolumeSource.DownwardAPI.DefaultMode

			for _, val := range val.VolumeSource.DownwardAPI.Items {
				var item = api.DownwardAPIVolumeFile{}

				item.Mode = val.Mode
				item.Path = val.Path

				if val.FieldRef != nil {
					item.FieldRef = new(api.ObjectFieldSelector)
					item.FieldRef.APIVersion = val.FieldRef.APIVersion
					item.FieldRef.FieldPath = val.FieldRef.FieldPath
				}

				if val.ResourceFieldRef != nil {
					item.ResourceFieldRef = new(api.ResourceFieldSelector)
					item.ResourceFieldRef.ContainerName = val.ResourceFieldRef.ContainerName
					item.ResourceFieldRef.Divisor = val.ResourceFieldRef.Divisor
					item.ResourceFieldRef.Resource = val.ResourceFieldRef.Resource
				}

				volume.VolumeSource.DownwardAPI.Items = append(volume.VolumeSource.DownwardAPI.Items, item)
			}

		}

		if val.VolumeSource.FC != nil {
			volume.VolumeSource.FC = new(api.FCVolumeSource)
			volume.VolumeSource.FC.FSType = val.VolumeSource.FC.FSType
			volume.VolumeSource.FC.Lun = val.VolumeSource.FC.Lun
			volume.VolumeSource.FC.ReadOnly = val.VolumeSource.FC.ReadOnly
			volume.VolumeSource.FC.TargetWWNs = val.VolumeSource.FC.TargetWWNs
		}

		if val.VolumeSource.AzureFile != nil {
			volume.VolumeSource.AzureFile = new(api.AzureFileVolumeSource)
			volume.VolumeSource.AzureFile.ReadOnly = val.VolumeSource.AzureFile.ReadOnly
			volume.VolumeSource.AzureFile.SecretName = val.VolumeSource.AzureFile.SecretName
			volume.VolumeSource.AzureFile.ShareName = val.VolumeSource.AzureFile.ShareName
		}

		if val.VolumeSource.ConfigMap != nil {
			volume.VolumeSource.ConfigMap = new(api.ConfigMapVolumeSource)
			volume.VolumeSource.ConfigMap.DefaultMode = val.VolumeSource.ConfigMap.DefaultMode
			volume.VolumeSource.ConfigMap.Name = val.VolumeSource.ConfigMap.Name
			volume.VolumeSource.ConfigMap.LocalObjectReference = api.LocalObjectReference{
				Name: val.VolumeSource.ConfigMap.LocalObjectReference.Name,
			}

			for _, val := range val.VolumeSource.ConfigMap.Items {
				var item = api.KeyToPath{}

				item.Key = val.Key
				item.Mode = val.Mode
				item.Path = val.Path

				volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, item)
			}
		}

		if val.VolumeSource.VsphereVolume != nil {
			volume.VolumeSource.VsphereVolume = new(api.VsphereVirtualDiskVolumeSource)
			volume.VolumeSource.VsphereVolume.FSType = val.VolumeSource.VsphereVolume.FSType
			volume.VolumeSource.VsphereVolume.VolumePath = val.VolumeSource.VsphereVolume.VolumePath
		}

		if val.VolumeSource.AzureDisk != nil {
			volume.VolumeSource.AzureDisk = new(api.AzureDiskVolumeSource)
			volume.VolumeSource.AzureDisk.FSType = val.VolumeSource.AzureDisk.FSType
			volume.VolumeSource.AzureDisk.ReadOnly = val.VolumeSource.AzureDisk.ReadOnly

			if val.VolumeSource.AzureDisk.CachingMode != nil {
				var item = api.AzureDataDiskCachingMode(*val.VolumeSource.AzureDisk.CachingMode)
				volume.VolumeSource.AzureDisk.CachingMode = &item
			}

			volume.VolumeSource.AzureDisk.DataDiskURI = val.VolumeSource.AzureDisk.DataDiskURI
			volume.VolumeSource.AzureDisk.DiskName = val.VolumeSource.AzureDisk.DiskName
		}

		out.Volumes = append(out.Volumes, volume)
	}

	return nil
}

func Convert_v1_Container_to_api_Container(in *v1.Container, out *api.Container) error {

	if in == nil {
		return errors.New("Error: incoming data can not be nil")
	}

	out.Name = in.Name
	out.Image = in.Image
	out.Command = in.Command
	out.Args = in.Args
	out.WorkingDir = in.WorkingDir

	for _, val := range in.Ports {
		var port = api.ContainerPort{}

		port.Name = val.Name
		port.HostPort = val.HostPort
		port.ContainerPort = val.ContainerPort
		port.Protocol = api.Protocol(val.Protocol)
		port.HostIP = val.HostIP

		out.Ports = append(out.Ports, port)
	}

	for _, val := range in.Env {
		var env = api.EnvVar{}

		env.Name = val.Name
		env.Value = val.Value
		out.Env = append(out.Env, env)

		if val.ValueFrom != nil {
			env.ValueFrom = new(api.EnvVarSource)

			if val.ValueFrom.FieldRef != nil {
				env.ValueFrom.FieldRef = new(api.ObjectFieldSelector)
				env.ValueFrom.FieldRef.APIVersion = val.ValueFrom.FieldRef.APIVersion
				env.ValueFrom.FieldRef.FieldPath = val.ValueFrom.FieldRef.FieldPath
			}

			if val.ValueFrom.ResourceFieldRef != nil {
				env.ValueFrom.ResourceFieldRef = new(api.ResourceFieldSelector)
				env.ValueFrom.ResourceFieldRef.ContainerName = val.ValueFrom.ResourceFieldRef.ContainerName
				env.ValueFrom.ResourceFieldRef.Divisor = val.ValueFrom.ResourceFieldRef.Divisor
				env.ValueFrom.ResourceFieldRef.Resource = val.ValueFrom.ResourceFieldRef.Resource
			}

			if val.ValueFrom.ConfigMapKeyRef != nil {
				env.ValueFrom.ConfigMapKeyRef = new(api.ConfigMapKeySelector)
				env.ValueFrom.ConfigMapKeyRef.Key = val.ValueFrom.ConfigMapKeyRef.Key
				env.ValueFrom.ConfigMapKeyRef.Name = val.ValueFrom.ConfigMapKeyRef.Name
				env.ValueFrom.ConfigMapKeyRef.LocalObjectReference.Name = val.ValueFrom.ConfigMapKeyRef.LocalObjectReference.Name
			}

			if val.ValueFrom.SecretKeyRef != nil {
				env.ValueFrom.SecretKeyRef = new(api.SecretKeySelector)
				env.ValueFrom.SecretKeyRef.Name = val.ValueFrom.SecretKeyRef.Name
				env.ValueFrom.SecretKeyRef.Key = val.ValueFrom.SecretKeyRef.Key
				env.ValueFrom.SecretKeyRef.LocalObjectReference.Name = val.ValueFrom.SecretKeyRef.LocalObjectReference.Name
			}
		}
	}

	out.Resources.Limits = api.ResourceList{}
	for key, val := range in.Resources.Limits {
		var item = resource.Quantity{}
		item.Format = val.Format

		out.Resources.Limits[api.ResourceName(key)] = item
	}

	out.Resources.Requests = api.ResourceList{}
	for key, val := range in.Resources.Requests {
		out.Resources.Requests[api.ResourceName(key)] = val
	}

	for _, val := range in.VolumeMounts {
		out.VolumeMounts = append(out.VolumeMounts, api.VolumeMount{
			Name:      val.Name,
			ReadOnly:  val.ReadOnly,
			MountPath: val.MountPath,
			SubPath:   val.SubPath,
		})
	}

	if in.LivenessProbe != nil {
		out.LivenessProbe = new(api.Probe)

		if in.LivenessProbe.Exec != nil {
			out.LivenessProbe.Exec = new(api.ExecAction)
			out.LivenessProbe.Exec.Command = in.LivenessProbe.Exec.Command
		}

		if in.LivenessProbe.HTTPGet != nil {
			out.LivenessProbe.HTTPGet = new(api.HTTPGetAction)
			out.LivenessProbe.HTTPGet.Host = in.LivenessProbe.HTTPGet.Host
			out.LivenessProbe.HTTPGet.Path = in.LivenessProbe.HTTPGet.Path
			out.LivenessProbe.HTTPGet.Port = in.LivenessProbe.HTTPGet.Port
			out.LivenessProbe.HTTPGet.Scheme = api.URIScheme(in.LivenessProbe.HTTPGet.Scheme)

			for _, val := range in.LivenessProbe.HTTPGet.HTTPHeaders {
				out.LivenessProbe.HTTPGet.HTTPHeaders = append(out.LivenessProbe.HTTPGet.HTTPHeaders, api.HTTPHeader{
					Name:  val.Name,
					Value: val.Value,
				})
			}
		}

		if in.LivenessProbe.TCPSocket != nil {
			out.LivenessProbe.TCPSocket = new(api.TCPSocketAction)
			out.LivenessProbe.TCPSocket.Port = in.LivenessProbe.TCPSocket.Port
		}

		out.LivenessProbe.FailureThreshold = in.LivenessProbe.FailureThreshold
		out.LivenessProbe.InitialDelaySeconds = in.LivenessProbe.InitialDelaySeconds
		out.LivenessProbe.PeriodSeconds = in.LivenessProbe.PeriodSeconds
		out.LivenessProbe.SuccessThreshold = in.LivenessProbe.SuccessThreshold
		out.LivenessProbe.TimeoutSeconds = in.LivenessProbe.TimeoutSeconds
	}

	if in.ReadinessProbe != nil {
		out.ReadinessProbe = new(api.Probe)

		if in.ReadinessProbe.Exec != nil {
			out.ReadinessProbe.Exec = new(api.ExecAction)
			out.ReadinessProbe.Exec.Command = in.ReadinessProbe.Exec.Command
		}

		if in.ReadinessProbe.HTTPGet != nil {
			out.ReadinessProbe.HTTPGet = new(api.HTTPGetAction)
			out.ReadinessProbe.HTTPGet.Host = in.ReadinessProbe.HTTPGet.Host
			out.ReadinessProbe.HTTPGet.Path = in.ReadinessProbe.HTTPGet.Path
			out.ReadinessProbe.HTTPGet.Port = in.ReadinessProbe.HTTPGet.Port
			out.ReadinessProbe.HTTPGet.Scheme = api.URIScheme(in.ReadinessProbe.HTTPGet.Scheme)

			for _, val := range in.ReadinessProbe.HTTPGet.HTTPHeaders {
				out.ReadinessProbe.HTTPGet.HTTPHeaders = append(out.ReadinessProbe.HTTPGet.HTTPHeaders, api.HTTPHeader{
					Name:  val.Name,
					Value: val.Value,
				})
			}
		}

		if in.ReadinessProbe.TCPSocket != nil {
			out.ReadinessProbe.TCPSocket = new(api.TCPSocketAction)
			out.ReadinessProbe.TCPSocket.Port = in.ReadinessProbe.TCPSocket.Port
		}

		out.ReadinessProbe.FailureThreshold = in.ReadinessProbe.FailureThreshold
		out.ReadinessProbe.InitialDelaySeconds = in.ReadinessProbe.InitialDelaySeconds
		out.ReadinessProbe.PeriodSeconds = in.ReadinessProbe.PeriodSeconds
		out.ReadinessProbe.SuccessThreshold = in.ReadinessProbe.SuccessThreshold
		out.ReadinessProbe.TimeoutSeconds = in.ReadinessProbe.TimeoutSeconds
	}

	if in.Lifecycle != nil {
		out.Lifecycle = new(api.Lifecycle)

		if in.Lifecycle.PostStart != nil {
			out.Lifecycle.PostStart = new(api.Handler)

			if in.Lifecycle.PostStart.Exec != nil {
				out.Lifecycle.PostStart.Exec = new(api.ExecAction)
				out.Lifecycle.PostStart.Exec.Command = in.Lifecycle.PostStart.Exec.Command
			}

			for _, val := range in.Lifecycle.PostStart.HTTPGet.HTTPHeaders {
				out.Lifecycle.PostStart.HTTPGet.HTTPHeaders = append(out.Lifecycle.PostStart.HTTPGet.HTTPHeaders, api.HTTPHeader{
					Name:  val.Name,
					Value: val.Value,
				})
			}

			if in.Lifecycle.PostStart.TCPSocket != nil {
				out.Lifecycle.PostStart.TCPSocket = new(api.TCPSocketAction)
				out.Lifecycle.PostStart.TCPSocket.Port = in.Lifecycle.PostStart.TCPSocket.Port
			}

		}

		if in.Lifecycle.PreStop != nil {
			out.Lifecycle.PreStop = new(api.Handler)

			if in.Lifecycle.PreStop.Exec != nil {
				out.Lifecycle.PreStop.Exec = new(api.ExecAction)
				out.Lifecycle.PreStop.Exec.Command = in.Lifecycle.PreStop.Exec.Command
			}

			for _, val := range in.Lifecycle.PreStop.HTTPGet.HTTPHeaders {
				out.Lifecycle.PreStop.HTTPGet.HTTPHeaders = append(out.Lifecycle.PreStop.HTTPGet.HTTPHeaders, api.HTTPHeader{
					Name:  val.Name,
					Value: val.Value,
				})
			}

			if in.Lifecycle.PreStop.TCPSocket != nil {
				out.Lifecycle.PreStop.TCPSocket = new(api.TCPSocketAction)
				out.Lifecycle.PreStop.TCPSocket.Port = in.Lifecycle.PreStop.TCPSocket.Port
			}
		}
	}

	out.TerminationMessagePath = in.TerminationMessagePath

	out.ImagePullPolicy = api.PullPolicy(in.ImagePullPolicy)

	if in.SecurityContext != nil {

		if in.SecurityContext.Capabilities != nil {
			out.SecurityContext.Capabilities = new(api.Capabilities)

			for _, val := range in.SecurityContext.Capabilities.Add {
				out.SecurityContext.Capabilities.Add = append(out.SecurityContext.Capabilities.Add, api.Capability(val))
			}

			for _, val := range in.SecurityContext.Capabilities.Drop {
				out.SecurityContext.Capabilities.Drop = append(out.SecurityContext.Capabilities.Drop, api.Capability(val))
			}
		}

		out.SecurityContext = new(api.SecurityContext)
		out.SecurityContext.ReadOnlyRootFilesystem = in.SecurityContext.Privileged
		out.SecurityContext.RunAsNonRoot = in.SecurityContext.RunAsNonRoot
		out.SecurityContext.RunAsUser = in.SecurityContext.RunAsUser

		if in.SecurityContext.SELinuxOptions != nil {
			out.SecurityContext.SELinuxOptions = new(api.SELinuxOptions)
			in.SecurityContext.SELinuxOptions.Level = in.SecurityContext.SELinuxOptions.Level
			in.SecurityContext.SELinuxOptions.Role = in.SecurityContext.SELinuxOptions.Role
			in.SecurityContext.SELinuxOptions.Type = in.SecurityContext.SELinuxOptions.Type
			in.SecurityContext.SELinuxOptions.User = in.SecurityContext.SELinuxOptions.User
		}

		out.SecurityContext.Privileged = in.SecurityContext.Privileged
	}

	out.Stdin = in.Stdin
	out.StdinOnce = in.StdinOnce
	out.TTY = in.TTY

	return nil
}

func Convert_v1_OwnerReference_To_api_OwnerReference(in *v1.OwnerReference, out *api.OwnerReference) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.APIVersion = in.APIVersion
	out.Kind = in.Kind
	out.Name = in.Name
	out.UID = types.UID(in.UID)
	out.Controller = in.Controller

	return nil
}

func Convert_v1beta1_Deployment_To_extensions_Deployment(in *v1beta1.Deployment, out *extensions.Deployment) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	Set_defaults_extensions_deployment(out)

	if err := Convert_unversioned_TypeMeta_to_unversioned_TypeMeta(&in.TypeMeta, &out.TypeMeta); err != nil {
		return err
	}

	if err := Convert_v1_ObjectMeta_to_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta); err != nil {
		return err
	}

	if err := Convert_v1beta1_DeploymentSpec_to_extensions_DeploymentSpec(&in.Spec, &out.Spec); err != nil {
		return err
	}

	if err := Convert_v1beta1_DeploymentStatus_to_extensions_DeploymentStatus(&in.Status, &out.Status); err != nil {
		return err
	}

	return nil
}

func Convert_v1beta1_DeploymentSpec_to_extensions_DeploymentSpec(in *v1beta1.DeploymentSpec, out *extensions.DeploymentSpec) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	if in.Replicas != nil {
		out.Replicas = *in.Replicas
	}

	if in.Selector != nil {
		out.Selector = new(unversioned.LabelSelector)
		if err := Convert_v1beta1_LabelSelector_to_unversioned_LabelSelector(in.Selector, out.Selector); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}

	if err := Convert_v1_PodTemplateSpec_to_api_PodTemplateSpec(&in.Template, &out.Template); err != nil {
		return err
	}

	if err := Convert_v1beta1_DeploymentStrategy_to_extensions_DeploymentStrategy(&in.Strategy, &out.Strategy); err != nil {
		return err
	}

	out.RevisionHistoryLimit = in.RevisionHistoryLimit
	out.MinReadySeconds = in.MinReadySeconds
	out.Paused = in.Paused
	if in.RollbackTo != nil {
		out.RollbackTo = new(extensions.RollbackConfig)
		out.RollbackTo.Revision = in.RollbackTo.Revision
	} else {
		out.RollbackTo = nil
	}

	return nil
}

func Convert_v1beta1_DeploymentStrategy_to_extensions_DeploymentStrategy(in *v1beta1.DeploymentStrategy, out *extensions.DeploymentStrategy) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.Type = extensions.DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(extensions.RollingUpdateDeployment)
		if err := Convert_v1beta1_RollingUpdateDeployment_to_extensions_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}

	return nil
}

func Convert_v1beta1_RollingUpdateDeployment_to_extensions_RollingUpdateDeployment(in *v1beta1.RollingUpdateDeployment, out *extensions.RollingUpdateDeployment) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.MaxSurge = intstr.IntOrString{
		Type:   in.MaxSurge.Type,
		IntVal: in.MaxSurge.IntVal,
		StrVal: in.MaxSurge.StrVal,
	}

	out.MaxUnavailable = intstr.IntOrString{
		Type:   in.MaxUnavailable.Type,
		IntVal: in.MaxUnavailable.IntVal,
		StrVal: in.MaxUnavailable.StrVal,
	}

	return nil
}

func Convert_v1beta1_LabelSelector_to_unversioned_LabelSelector(in *v1beta1.LabelSelector, out *unversioned.LabelSelector) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.MatchLabels = in.MatchLabels
	if in.MatchExpressions != nil {
		in, out := &in.MatchExpressions, &out.MatchExpressions
		*out = make([]unversioned.LabelSelectorRequirement, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_LabelSelectorRequirement_to_unversioned_LabelSelectorRequirement(&(*in)[i], &(*out)[i]); err != nil {
				return err
			}
		}
	} else {
		out.MatchExpressions = nil
	}

	return nil
}

func Convert_v1beta1_LabelSelectorRequirement_to_unversioned_LabelSelectorRequirement(in *v1beta1.LabelSelectorRequirement, out *unversioned.LabelSelectorRequirement) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.Key = in.Key
	out.Operator = unversioned.LabelSelectorOperator(in.Operator)
	out.Values = in.Values

	return nil
}

func Convert_v1beta1_DeploymentStatus_to_extensions_DeploymentStatus(in *v1beta1.DeploymentStatus, out *extensions.DeploymentStatus) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.ObservedGeneration = in.ObservedGeneration
	out.Replicas = in.Replicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.AvailableReplicas = in.AvailableReplicas
	out.UnavailableReplicas = in.UnavailableReplicas

	return nil
}

func Convert_unversioned_TypeMeta_to_unversioned_TypeMeta(in, out *unversioned.TypeMeta) error {

	if in == nil {
		return error_incoming_data
	}

	if out == nil {
		return error_outcoming_data
	}

	out.APIVersion = in.APIVersion
	out.Kind = in.Kind

	return nil
}
