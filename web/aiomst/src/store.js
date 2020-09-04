import axios from 'axios';
import Vue from 'vue';
import Vuex from 'vuex';
import {Howl} from 'howler';

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    gResp: null,
    playback: null,
    playlist: null,
  },
  getters: {
    currentArtist: state => {
      return state.gResp["artists"].length === 1 ? state.gResp : null
    },      
    getArtists: state => {
      return state.gResp["artists"];
    },

    currentAlbum: state => {
      return state.gResp["albums"].length === 1 ? state.gResp : null
    },      
    getAlbums: state => {
      return state.gResp["albums"];
    },

    currentSong: state => {
      return state.playback
    },
    getSongs: state => {
      return state.gResp["songs"];
    },
    getPlaylist: state => {
      return state.playlist;
    }
  },
  mutations: {
    changeGoResp: (state, payload) => {
      state.gResp = payload.gResp;
      if (payload.gResp["songs"].length > 0) {
        state.playlist = payload.gResp["songs"]
      }
    },
    setSongState: (state, payload) => {
      state.playback = payload.playback
    },
    setPlaylist: (state, payload) => {
      state.playlist = payload
    }
  },
  actions: {
    makeApiRequest: (context, navInfo) => {
      return new Promise((resolve, reject) => {
        axios
          .get(navInfo.path, {
            baseURL: "http://127.0.0.1:8090",
            params: navInfo.params,
          })
          .then((x) => {
            context.commit('changeGoResp', {
              gResp: x.data,
            });
            resolve();
          })
          .catch((x) => {
            console.log(x);
            reject();
          });
      });
    },
    streamAudio: (context, song) => {
      // Check if a request for the current playing song has been made
      // If so, ignore it
      var playing = context.getters.currentSong;
      if (playing && playing.song.id === song.id) {
        return
      }
      
      return new Promise((resolve, reject) => {
        axios
          .get("http://localhost:8090/stream?id=" + song.id, {
            headers: {
              "Accept": "text/event-stream",
            },
            responseType: "arraybuffer"
          })
          .then(response => {
            // If there is a song playing, stop it
            if (context.getters.currentSong) {
              var playback = context.getters.currentSong
              playback.howl.unload()
            }
                  
            var blob = new Blob(
              [response.data], 
              {type: response.headers["content-type"]}
            );

            var url = URL.createObjectURL(blob);
            var howlInfo = new Howl({
              src: [url],
              html5: true,
              volume: .5,
              preload: false,
              onend: () => {
                var list = context.getters.getPlaylist
                var indexCurr = list.indexOf(context.getters.currentSong.song)
                var next = list[0]

                if (indexCurr + 1 < list.length) {
                  // add playlist loop condition here possibly
                  next = list[indexCurr + 1]
                }
                
                context.dispatch('streamAudio', next)
              }
            });
      
            var howlId = howlInfo.play();            
            context.commit(
              'setSongState', {
                'playback': {
                  'howl': howlInfo,
                  'howlId': howlId,
                  'song': song,
                }
              }
            )
            resolve();
          })
          .catch(error => {
            console.error(error);
            reject();
          });
      }); 
    }
  }
});