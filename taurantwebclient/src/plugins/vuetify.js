import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';

Vue.use(Vuetify);

export default new Vuetify({
    theme:{
        themes: {
            dark: {
              primary: '#4b0082'
            }
        },
        dark: true
    }
})
