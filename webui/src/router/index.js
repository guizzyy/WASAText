import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "../views/LoginView.vue";
import ConversationView from "../views/ConversationView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/conversations', component: HomeView},
		{path: '/conversations/:convID', component: ConversationView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
