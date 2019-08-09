import {Config} from '../../shared/model/Config';

export function fakeAppConfig(): Config {
  return {
    default_node_count: 3,
    show_demo_info: false,
    show_terms_of_service: false,
    share_kubeconfig: false,
    openstack: {
      wizard_use_default_user: false,
    },
    cleanup_cluster: false,
    enforce_cleanup_cluster: false,
  };
}
