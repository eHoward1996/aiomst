<template>
  <cmpntLoadBar v-if="!finishedLoading"></cmpntLoadBar>
  <v-container v-else fluid>
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
      <eleSongsList></eleSongsList>
    </v-container>
    <v-container v-if="checkNoResults">
      <h1>No Results Found</h1>
    </v-container>
  </v-container>
</template>

<script>
import { mapGetters, mapState } from 'vuex';
import cmpntLoadBar from '@/components/layout/cmpntLoadBar.vue';
import eleCardList  from '@/components/elements/eleCardList.vue';
import eleSongsList from '@/components/elements/eleSongsList.vue';

export default {
  name: 'Search',
  components: {cmpntLoadBar, eleCardList, eleSongsList},
  data: function() {
    return {
      finishedLoading: false,
    }
  },
  computed: {
    ...mapGetters({
      artists: 'getArtists',
      albums:  'getAlbums',
      songs:   'getSongs',
    }),
    ...mapState({
      apiResp: 'gResp', 
    }),
    checkArtistsInState: function() {
      if (!this.artists) {
        return false
      }
      return this.artists.length > 0;
    },
    checkAlbumsInState: function() {
      if (!this.albums) {
        return false
      }
      return this.albums.length > 0;
    },
    checkSongsInState: function() {
      if (!this.songs) {
        return false
      }
      return this.songs.length > 0;
    },
    checkNoResults: function() {
      return this.artists.length === 0 &&
             this.albums.length === 0 && 
             this.songs.length === 0
    }
  },
  watch: {
    apiResp: function() {
      this.finishedLoading = true
    }
  }
}
</script>