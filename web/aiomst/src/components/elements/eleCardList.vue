<template v-if="hasObjects">
  <v-container fluid>
    <v-row 
      v-for="i in rowCount"
      :key="i">
      <v-col 
        cols="3"
        v-for="obj in objs.slice((i - 1) * 4, i * 4)"
        :key="obj.id">
        <v-card 
          width="300"
          height="360"
          max-width="300"
          max-height="360"
          v-on:click="objAddr(obj.id)">
          <v-img
            height="300"
            width="300"
            :src="artSrc(obj)">
          </v-img>
          <v-tooltip bottom>
            <template v-slot:activator="{ on, attrs }">
              <v-card-title
                v-bind="attrs"
                v-on="on"
                class="text-no-wrap">
                {{getShortTitle(obj)}}
              </v-card-title>
            </template>
            <span>{{obj.title}}</span>
          </v-tooltip>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'eleCardList',
  props: ['req'],
  computed: {
    objs() {
      switch (this.req) {
        case 'albums':
          return this.$store.getters.albums;
        case 'artists':
          return this.$store.getters.artists;
        default:
          return [];
      }
    },
    hasObjects() {
      return this.objs.length > 0;
    },
    rowCount() {     
      return Math.ceil(this.objs.length / 4);
    },
  },
  methods: {
    artSrc: function(obj) {
      if (obj.artId !== 0) {
        return "http://localhost:8090/art?id=" + obj.artId;
      }
    },
    objAddr: function(objId) {
      let p = '';
      switch (this.req) {
        case 'albums':
          p = 'album';
          break;
        case 'artists':
          p = 'artist';
          break;
      }
      this.$router.push({
        path: p, 
        query: {id: objId},
      })
    },
    getShortTitle: function(obj) {
      if (obj.title.length > 25) {
        return obj.title.slice(0, 23) + "..."
      }
      return obj.title
    },
  },
}
</script>