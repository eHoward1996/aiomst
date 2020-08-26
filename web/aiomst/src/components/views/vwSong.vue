<template>
  <cmpntLoadBar v-if="!finishedLoading"></cmpntLoadBar>
  <v-container v-else>
    <eleSongList></eleSongList>
  </v-container>
</template>

<script>
import eleSongList from '@/components/elements/eleSongList.vue';
import cmpntLoadBar from '@/components/layout/cmpntLoadBar.vue';

export default {
  name: 'Song',
  components: {eleSongList, cmpntLoadBar},
  created: function() {
    this.$store.dispatch('makeApiRequest', {path: '/songs'}).then(() => {
      this.finishedLoading = true;
    });
  },
  data: function() {
    return {
      finishedLoading: false,
    }
  },
}
</script>