import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    gResp: null,
  },
  getters: {
    getGoResp: state => {
      return state.gResp
    }
  },
  mutations: {
    changeGoResp: (state, payload) => {
      state.gResp = payload.gResp;
    }
  },
  actions: {

  }
});

// store.watch((state) => {
//   console.log(state.gResp)
// });