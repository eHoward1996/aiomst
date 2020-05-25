import Vue from 'vue';
import VueRouter from 'vue-router';

import Index   from './components/vwIndex';
import Albums  from './components/vwAlbums';
import Artists from './components/vwArtists';
import Songs   from './components/vwSongs';
import Search  from './components/vwSearch';

Vue.use(VueRouter);
export default new VueRouter({
 mode: 'history',
 routes: [
   {path: '/',        component: Index},
   {path: '/albums',  component: Albums},
   {path: '/artists', component: Artists},
   {path: '/songs',   component: Songs},
   {path: '/search',  component: Search}
 ]
});