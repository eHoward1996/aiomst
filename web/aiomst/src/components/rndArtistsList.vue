<template v-if="gResp['artists']">
  <div>
    <h1>Artists <v-icon color="primary" x-large>mdi-account</v-icon></h1>
    <v-row
      v-for="i in rowCount"
      :key="i">
      <v-col
        v-for="artist in gResp['artists'].slice((i - 1) * 3, i * 3)"
        :key="artist.id">
        <v-card 
          max-width="300"
          max-height="450">
          <v-img
            height="300"
            width="300"
            :src="artSrc(artist)">
          </v-img>
          <v-card-title>{{artist.title}}</v-card-title>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  name: 'rndArtistsList',
  computed: {
    gResp() {
      return this.$store.state.gResp;
    },
    rowCount() {     
      return Math.ceil(this.gResp['albums'].length / 3);
    },
  },
  methods: {
    artSrc: function(artist) {
      return "http://localhost:8090/art/" + artist.artId;
    }
  },
}
</script>