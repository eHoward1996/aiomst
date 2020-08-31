<template>
  <v-app>
    <cmpntAppBar></cmpntAppBar>
    <v-content class="main">
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-content>
    <elePlayerBar></elePlayerBar>
  </v-app>
</template>

<script>
import cmpntAppBar from '@/components/layout/cmpntAppBar.vue';
import elePlayerBar from '@/components/player/currentTrack.vue';

export default {
  name: 'App',
  components: {cmpntAppBar, elePlayerBar},
}
</script>

<style scoped lang="sass">
  .main
    margin-bottom: 5%
</style>