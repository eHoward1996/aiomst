import Vue from 'vue';
import App from './App.vue';
import axios from 'axios';
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
    '$route'(to)  {
      axios
        .get(to.path, {
          baseURL: "http://127.0.0.1:8090",
          params: to.query,
        })
        .then((x) => {
          this.$store.commit('changeGoResp', {
            gResp: x.data,
          });
        })
        .catch((x) => {
          this.gResp = x;
        });
    }
  }
}).$mount('#app')
