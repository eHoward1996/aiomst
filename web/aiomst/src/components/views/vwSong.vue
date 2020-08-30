<template>
  <cmpntLoadBar v-if="!finishedLoading"></cmpntLoadBar>
  <v-container v-else>
    <eleSongsList></eleSongsList>
  </v-container>
</template>

<script>
import eleSongsList from '@/components/elements/eleSongsList.vue';
import cmpntLoadBar from '@/components/layout/cmpntLoadBar.vue';

export default {
  name: 'Song',
  components: {eleSongsList, cmpntLoadBar},
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