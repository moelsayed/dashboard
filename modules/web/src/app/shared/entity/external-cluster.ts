// Copyright 2020 The Kubermatic Kubernetes Platform contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {KubeOneClusterSpec} from '@shared/entity/kubeone-cluster';
import {StatusIcon} from '@shared/utils/health-status';
import _ from 'lodash';
import {BringYourOwnCloudSpec, ContainerRuntime} from './cluster';
import {AKSCloudSpec, AKSClusterSpec} from './provider/aks';
import {EKSCloudSpec, EKSClusterSpec} from './provider/eks';
import {GKECloudSpec, GKEClusterSpec} from './provider/gke';

export enum DeleteExternalClusterAction {
  Delete = 'delete',
  Disconnect = 'disconnect',
}

export enum ExternalClusterProvider {
  Custom = 'bringyourown',
  AKS = 'aks',
  EKS = 'eks',
  GKE = 'gke',
  KubeOne = 'kubeone',
}

const PROVIDER_DISPLAY_NAMES = new Map<ExternalClusterProvider, string>([
  [ExternalClusterProvider.AKS, 'AKS'],
  [ExternalClusterProvider.EKS, 'EKS'],
  [ExternalClusterProvider.GKE, 'GKE'],
  [ExternalClusterProvider.Custom, 'Custom'],
  [ExternalClusterProvider.KubeOne, 'KubeOne'],
]);

export function getExternalProviderDisplayName(provider: ExternalClusterProvider): string {
  return PROVIDER_DISPLAY_NAMES.get(provider);
}

export class ExternalCluster {
  creationTimestamp?: Date;
  deletionTimestamp?: Date;
  id?: string;
  name: string;
  labels?: object;
  spec: ExternalClusterSpec;
  cloud: ExternalCloudSpec;
  status: ExternalClusterStatus;

  static newEmptyKubeOneClusterEntity(): ExternalCluster {
    return {
      name: '',
      cloud: {
        kubeOne: KubeOneClusterSpec.newEmptyClusterEntity(),
      } as ExternalCloudSpec,
    } as ExternalCluster;
  }

  static getProvider(cloud: ExternalCloudSpec): ExternalClusterProvider {
    if (!cloud) {
      return null;
    }
    const providers = Object.keys(cloud);
    return providers.pop()?.toLowerCase() as ExternalClusterProvider;
  }

  static getProviderDisplayName(cloud: ExternalCloudSpec): string {
    return getExternalProviderDisplayName(ExternalCluster.getProvider(cloud));
  }

  static isDeleted(cluster: ExternalCluster): boolean {
    if (!cluster) {
      return false;
    }
    return (
      !!cluster.deletionTimestamp ||
      cluster.status?.state === ExternalClusterState.Deleting ||
      (cluster.status?.state === ExternalClusterState.Error &&
        (cluster.status?.statusMessage?.includes('NotFound') || cluster.status?.statusMessage?.includes('Not Found')))
    );
  }

  static getStatusMessage(cluster: ExternalCluster): string {
    const state = _.capitalize(cluster?.status?.state);
    return cluster?.status?.statusMessage ? `${state}: ${cluster?.status.statusMessage}` : state;
  }

  static getStatusIcon(cluster: ExternalCluster): string {
    if (cluster?.status?.state) {
      return ExternalClusterStateIcon[cluster.status.state] || StatusIcon.Unknown;
    }
    return StatusIcon.Unknown;
  }
}

export class ExternalClusterSpec {
  aksclusterSpec?: AKSClusterSpec;
  eksclusterSpec?: EKSClusterSpec;
  gkeclusterSpec?: GKEClusterSpec;
  version?: string;
  containerRuntime?: ContainerRuntime;
}

export class ExternalCloudSpec {
  aks?: AKSCloudSpec;
  eks?: EKSCloudSpec;
  gke?: GKECloudSpec;
  bringYourOwn?: BringYourOwnCloudSpec;
  kubeOne?: KubeOneClusterSpec;
}

export enum ExternalClusterState {
  ConfigError = 'ConfigError',
  ConnectionError = 'ConnectionError',
  Deleting = 'Deleting',
  Error = 'Error',
  EtcdError = 'EtcdError',
  KubeClientError = 'KubeClientError',
  Provisioning = 'Provisioning',
  Reconciling = 'Reconciling',
  ReconcilingMigrate = 'ReconcilingMigrate',
  ReconcilingUpgrade = 'ReconcilingUpgrade',
  Running = 'Running',
  RuntimeError = 'RuntimeError',
  SSHError = 'SSHError',
  Starting = 'Starting',
  Stopped = 'Stopped',
  Stopping = 'Stopping',
  Unknown = 'Unknown',
  Warning = 'Warning',
}

const ExternalClusterStateIcon: Record<ExternalClusterState, StatusIcon> = {
  [ExternalClusterState.ConfigError]: StatusIcon.Error,
  [ExternalClusterState.ConnectionError]: StatusIcon.Error,
  [ExternalClusterState.Deleting]: StatusIcon.Error,
  [ExternalClusterState.Error]: StatusIcon.Error,
  [ExternalClusterState.EtcdError]: StatusIcon.Error,
  [ExternalClusterState.KubeClientError]: StatusIcon.Error,
  [ExternalClusterState.Provisioning]: StatusIcon.Pending,
  [ExternalClusterState.Reconciling]: StatusIcon.Pending,
  [ExternalClusterState.ReconcilingMigrate]: StatusIcon.Pending,
  [ExternalClusterState.ReconcilingUpgrade]: StatusIcon.Pending,
  [ExternalClusterState.Running]: StatusIcon.Running,
  [ExternalClusterState.RuntimeError]: StatusIcon.Error,
  [ExternalClusterState.SSHError]: StatusIcon.Error,
  [ExternalClusterState.Starting]: StatusIcon.Pending,
  [ExternalClusterState.Stopped]: StatusIcon.Stopped,
  [ExternalClusterState.Stopping]: StatusIcon.Pending,
  [ExternalClusterState.Unknown]: StatusIcon.Unknown,
  [ExternalClusterState.Warning]: StatusIcon.Warning,
};

export class ExternalClusterStatus {
  state: ExternalClusterState;
  statusMessage: string;
  // this property is used only for AKS Cluster.
  aks?: AKSClusterStatus;
}

class AKSClusterStatus {
  provisioningState: string;
  powerState: string;
}

export class ExternalClusterPatch {
  spec?: ExternalClusterSpecPatch;
}

export class ExternalClusterSpecPatch {
  version?: string;
}

export class ExternalClusterModel {
  name?: string;
  kubeconfig?: string;
  cloud?: ExternalCloudSpec;
  spec?: ExternalClusterSpec;

  static new(): ExternalClusterModel {
    return {name: ''};
  }
}
