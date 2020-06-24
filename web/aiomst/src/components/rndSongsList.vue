<template v-if="gResp['songs']">
  <div>
    <h1>Songs <v-icon color="primary" x-large>mdi-music-note</v-icon></h1>
    <v-list-item-group color="primary">
      <v-list-item
        v-for="song in gResp['songs']"
        :key="song.id">
        <v-list-item-avatar>
          {{song.track}}
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title v-html="song.title"></v-list-item-title>
            <v-list-item-subtitle v-html="song.artist"></v-list-item-subtitle>
        </v-list-item-content>
      <v-list-item-avatar>
        {{formatLength(song)}}
      </v-list-item-avatar>
      </v-list-item>
    </v-list-item-group>
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
  methods: {
    formatLength: function(song) {
      var t = parseInt(song.length);
      var minute = Math.floor(t / 60);
      var sec = t % 60;

      return minute + ":" + sec;
    }
  },
}
</script>