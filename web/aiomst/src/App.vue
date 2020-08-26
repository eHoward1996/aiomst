<template>
  <v-app>
    <cmpntAppBar></cmpntAppBar>
    <v-content>
      <v-container fluid>
        <router-view class='main'></router-view>
      </v-container>
      <elePlayerBar></elePlayerBar>
    </v-content>
  </v-app>
</template>

<script>
import cmpntAppBar from '@/components/layout/cmpntAppBar.vue';
import elePlayerBar from '@/components/player/elePlayerBar.vue';

export default {
  name: 'App',
  components: {cmpntAppBar, elePlayerBar},
}
</script>

<style scoped lang='sass'>
  .main
    position: relative
    display: flex
    flex-flow: column
    flex-grow: 1
    margin-bottom: 4%
    overflow-y: auto
    width: 100%
</style>