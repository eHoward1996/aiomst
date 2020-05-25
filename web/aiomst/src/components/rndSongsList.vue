<template v-if="gResp['songs']">
  <div>
    Songs <v-icon color="primary">mdi-music-note</v-icon>
    <v-card
      v-for="song in gResp['songs']"
      :key="song.title"
      max-width="300"
      max-height="350">
      <v-card-title>{{song.title}}</v-card-title>
    </v-card>
  </div>
</template>

<script>
export default {
  name: 'rndSongsList',
  computed: {
    gResp() {
      return this.$store.state.gResp;
    }
  },
}
</script>