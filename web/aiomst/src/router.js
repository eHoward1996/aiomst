import Vue from 'vue';
import VueRouter from 'vue-router';

import Index   from './components/views/vwIndex.vue';
import Album   from './components/views/vwAlbum.vue';
import Artist  from './components/views/vwArtist.vue';
import Song    from './components/views/vwSong.vue';
import Search  from './components/views/vwSearch.vue';

Vue.use(VueRouter);
export default new VueRouter({
  mode: 'history',
  routes: [
    {path: '/',        component: Index},
    {path: '/album',   component: Album},
    {path: '/artist',  component: Artist},
    {path: '/song',    component: Song},
    {path: '/search',  component: Search}
  ]
});