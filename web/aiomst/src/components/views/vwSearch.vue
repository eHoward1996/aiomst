<template>
  <v-container fluid>
    <v-container v-if="checkArtistsInState" fluid style="width: 75%;">
      <v-row><h1>Artists <v-icon color="primary" x-large>mdi-account</v-icon></h1></v-row>
      <eleCardList req="artists"></eleCardList>
    </v-container>
    <v-container v-if="checkAlbumsInState" fluid style="width: 75%;">
      <v-row><h1>Albums <v-icon color="primary" x-large>mdi-disc</v-icon></h1></v-row>
      <eleCardList req="albums"></eleCardList>
    </v-container>
    <v-container v-if="checkSongsInState" fluid style="width: 75%;">
      <v-row><h1>Songs <v-icon color="primary" x-large>mdi-music-note</v-icon></h1></v-row>
      <eleSongList></eleSongList>
    </v-container>
  </v-container>
</template>

<script>
import eleCardList from '@/components/elements/eleCardList.vue';
import eleSongList from '@/components/elements/eleSongList.vue';

export default {
  name: 'Search',
  components: {eleCardList, eleSongList},
  computed: {
    checkArtistsInState: function() {
      return this.artists;
    },
    checkAlbumsInState: function() {
      return this.albums;
    },
    checkSongsInState: function() {
      return this.songs;
    },
  },
  created: function() {
    let navInfo = {
      path: '/search',
      params: this.$route.query,
    };
    this.$store.dispatch('makeApiRequest', navInfo).then(() => {
      this.artists = this.$store.getters.artists;
      this.albums = this.$store.getters.albums;
      this.songs = this.$store.getters.songs;
    });
  },
  data: function() {
    return {
      artists: null,
      albums: null,
      songs: null
    }
  },
}
</script>