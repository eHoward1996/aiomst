<template v-if="gResp['artists']">
  <div>
    Artists <v-icon color="primary">mdi-account</v-icon>
    <v-card 
      v-for="artist in gResp['artists']"
      :key="artist.title"
      max-width="300"
      max-height="350">
      <v-card-title>{{artist.title}}</v-card-title>
    </v-card>
  </div>
</template>

<script>
export default {
  name: 'rndArtistsList',
  computed: {
    gResp() {
      return this.$store.state.gResp;
    }
  },
}
</script>