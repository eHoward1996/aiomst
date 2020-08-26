import Vue from 'vue';
import VueRouter from 'vue-router';

import Index   from './components/vwIndex';
import Album   from './components/vwAlbum';
import Artist  from './components/vwArtist';
import Song    from './components/vwSong';
import Search  from './components/vwSearch';

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