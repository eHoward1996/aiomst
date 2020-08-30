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
        <v-col cols="4"><eleSongsList></eleSongsList></v-col>
        <v-col cols="2"></v-col>
      </v-row>
      <v-row v-else>
        <eleCardList v-if="albums" req="albums"></eleCardList>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import eleCardList  from '@/components/elements/eleCardList.vue';
import eleSongsList from '@/components/elements/eleSongsList.vue';
import cmpntLoadBar from '@/components/layout/cmpntLoadBar.vue';

export default {
  name: 'Album',
  components: {eleCardList, eleSongsList, cmpntLoadBar},
  created: function() {
    let navInfo = {
      path: '/albums',
      params: this.$route.query,
    };

    this.$store.dispatch('makeApiRequest', navInfo).then(() => {
      this.albums = this.$store.getters.albums;
      this.isSingleAlbumView = false;
      if (this.albums.length === 1) {
        this.isSingleAlbumView = true;
      }
      this.finishedLoading = true;
    });
  },
  data: function() {
    return {
      albums: [],
      finishedLoading: false,
      isSingleAlbumView: false,
    }
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
    }
  },
  watch: {
    '$route.query.id': function() {
      this.finishedLoading = false;
      if (this.$route.query.id !== undefined) {
        this.albums = this.albums.filter(a => a.id === Number(this.$route.query.id));
        this.isSingleAlbumView = true;        
        this.finishedLoading = true;
        return
      }

      this.$store.dispatch('makeApiRequest', {path: '/albums'}).then(() => {
        this.albums = this.$store.getters.albums;
        this.isSingleAlbumView = false;
        this.finishedLoading = true;
      })
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