import Vue from 'vue';
import VueRouter from 'vue-router';

import Index   from './components/views/vwIndex';
import Album   from './components/views/vwAlbum';
import Artist  from './components/views/vwArtist';
import Song    from './components/views/vwSong';
import Search  from './components/views/vwSearch';

Vue.use(VueRouter);
export default new VueRouter({
  mode: 'history',
  routes: [
    {path: '/',        component: Index},
    {path: '/albums',   component: Album},
    {path: '/artists',  component: Artist},
    {path: '/songs',    component: Song},
    {path: '/search',  component: Search}
  ]
});