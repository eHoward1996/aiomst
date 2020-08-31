<template>
  <cmpntLoadBar v-if="!finishedLoading"></cmpntLoadBar>
  <div v-else>
    <v-container fluid style="width: 80%;">
      <v-row v-if="isSingleAlbumView">
        <v-col cols="2"></v-col>
        <v-col cols="4">
          <v-row><h1>{{this.albums[0].title}}</h1></v-row>
          <v-row>
            <h1 
              class="aName"
              @click="artistAddr()">
              {{this.albums[0].artist}}
            </h1>
          </v-row>
          <v-row>
            <v-img
              max-height="300"
              max-width="300"
              :src="artSrc()">
            </v-img>
          </v-row>
        </v-col>
        <v-col cols="4">
          <v-row>
            <v-col cols="6" class="text-left">{{getTracksLength()}}</v-col>
            <v-col cols="6" class="text-right">{{getPlayLength()}}</v-col>
          </v-row>
          <v-row>
            <eleSongsList></eleSongsList>
          </v-row>
        </v-col>
        <v-col cols="2"></v-col>
      </v-row>
   
      <v-row v-else>
        <eleCardList v-if="albums" req="albums"></eleCardList>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import {mapGetters, mapState} from 'vuex';
import eleCardList  from '@/components/elements/eleCardList.vue';
import eleSongsList from '@/components/elements/eleSongsList.vue';
import cmpntLoadBar from '@/components/layout/cmpntLoadBar.vue';

export default {
  name: 'Album',
  components: {eleCardList, eleSongsList, cmpntLoadBar},
  
  data: function() {
    return {
      finishedLoading: false,
      isSingleAlbumView: false,
    }
  },
  computed: {
    ...mapGetters({
      albums:       'getAlbums',
      currentAlbum: 'currentAlbum',
    }),
    ...mapState({
      apiResp: 'gResp', 
    })
  },
  methods: {
    artSrc: function() {
      if (this.albums[0].artId !== 0) {
        return "http://localhost:8090/art?id=" + this.albums[0].artId;
      }
    },
    artistAddr: function() {
      let aID = this.albums[0].artistId;
      this.$router.push({
        path: '/artist',
        query: {id: aID},
      });
    },
    getTracksLength: function() {
      return this.currentAlbum["songs"].length + " Tracks"
    },
    getPlayLength: function() {
      var sum = 0
      for (const s of this.currentAlbum["songs"]) {
        sum += s.length
      }

      var t = parseInt(sum);
      var minute = Math.floor(t / 60);

      const zeroPad = (num) => String(num).padStart(2, '0')
      var sec = zeroPad(t % 60);

      return minute + ":" + sec;
    }
  },
  watch: {
    apiResp: function() {
      this.isSingleAlbumView = false
      if (this.currentAlbum) {
        this.isSingleAlbumView = true
      }
      this.finishedLoading = true
    }
  }
}
</script>

<style scoped>
  .aName:hover {
    text-decoration: underline;
    cursor: pointer;
  }
</style>