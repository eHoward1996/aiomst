import Vue from 'vue';
import App from './App.vue';
import vuetify from './plugins/vuetify';
import router from './router.js';
import {store} from './store.js';

Vue.config.productionTip = false


new Vue({
  vuetify,
  render: h => h(App),
  router,
  store: store,
  watch: {
    '$route': (route) => {
      var apiPath = route.path;
      switch (apiPath) {
        case '/album':
          apiPath = '/albums';
          break;
        case '/artist':
          apiPath = '/artists';
          break;
        case '/search':
          break;
        case '/song':
          apiPath = '/songs';
          break;
        default:
          console.log('unknown path: ' + apiPath)
          break;
      }

      store.dispatch('makeApiRequest', {
        path: apiPath,
        params: route.query,
      });
    }
  }
}).$mount('#app')
