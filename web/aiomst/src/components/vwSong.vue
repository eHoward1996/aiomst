<template>
  <cmpntLoadBar v-if="!finishedLoading"></cmpntLoadBar>
  <v-container v-else>
    <rndSongsList></rndSongsList>
  </v-container>
</template>

<script>
import rndSongsList from '@/components/rndSongsList.vue';
import cmpntLoadBar from '@/components/cmpntLoadBar.vue';

export default {
  name: 'Song',
  components: {rndSongsList, cmpntLoadBar},
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