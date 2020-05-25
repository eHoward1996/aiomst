<template v-if="gResp['albums']">
  <div>
    Albums <v-icon color="primary">mdi-disc</v-icon>
    <v-card 
      v-for="album in gResp['albums']"
      :key="album.title"
      max-width="300"
      max-height="350">
      <v-card-title>{{album.title}}</v-card-title>
    </v-card>
  </div>
</template>

<script>
export default {
  name: 'rndAlbumsList',
  computed: {
    gResp() {
      return this.$store.state.gResp;
    }
  },
}
</script>