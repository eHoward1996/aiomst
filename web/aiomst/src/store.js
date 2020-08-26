import axios from 'axios';
import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    gResp: null,
  },
  getters: {
    artists: state => {
      return state.gResp["artists"];
    },
    albums: state => {
      return state.gResp["albums"];
    },
    songs: state => {
      return state.gResp["songs"];
    },
  },
  mutations: {
    changeGoResp: (state, payload) => {
      state.gResp = payload.gResp;
    }
  },
  actions: {
    makeApiRequest: (context, navInfo) => {
      return new Promise((resolve, reject) => {
        axios
          .get(navInfo.path, {
            baseURL: "http://127.0.0.1:8090",
            params: navInfo.params,
          })
          .then((x) => {
            context.commit('changeGoResp', {
              gResp: x.data,
            });
            resolve();
          })
          .catch((x) => {
            console.log(x);
            reject();
          });
      });
    }
  }
});