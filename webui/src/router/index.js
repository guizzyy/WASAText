import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "../views/LoginView.vue";
import ConversationView from "../views/ConversationView.vue";
import ProfileView from "../views/ProfileView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/conversations', component: HomeView},
		{path: '/conversations/:convID', component: ConversationView},
		{path: '/users/:uID/', component: ProfileView},
	]
})

export default router
