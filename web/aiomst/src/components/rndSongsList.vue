<template v-if="hasObjects">
  <v-container fluid>
    <v-list-item-group color="primary">
      <v-list-item
        v-for="song in songs"
        :key="song.id">
        <v-list-item-avatar>
          {{formatTrackNumber(song)}}
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
  </v-container>
</template>

<script>
export default {
  name: 'rndSongsList',
  computed: {
    songs() {
      return this.$store.getters.songs;
    },
    hasObjects() {
      return this.songs.length > 0;
    },
  },
  methods: {
    formatLength: function(song) {
      var t = parseInt(song.length);
      var minute = Math.floor(t / 60);

      const zeroPad = (num) => String(num).padStart(2, '0')
      var sec = zeroPad(t % 60);

      return minute + ":" + sec;
    },
    formatTrackNumber: function(song) {
      const zeroPad = (num) => String(num).padStart(2, '0')
      return zeroPad(song.track)
    }
  },
}
</script>