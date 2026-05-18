package com.freelanceplatform.android.ui.auth

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.freelanceplatform.android.data.local.TokenManager
import com.freelanceplatform.android.ui.home.HomeScreen
import com.freelanceplatform.android.ui.navigation.BottomNavBar
import com.freelanceplatform.android.ui.project.ProjectDetailScreen
import com.freelanceplatform.android.ui.project.ProjectListScreen
import com.freelanceplatform.android.ui.project.CreateProjectScreen
import com.freelanceplatform.android.ui.developer.DeveloperListScreen
import com.freelanceplatform.android.ui.developer.DeveloperDetailScreen
import com.freelanceplatform.android.ui.chat.ChatListScreen
import com.freelanceplatform.android.ui.chat.ChatDetailScreen
import com.freelanceplatform.android.ui.profile.ProfileScreen
import com.freelanceplatform.android.ui.profile.EditProfileScreen
import com.freelanceplatform.android.ui.wallet.WalletScreen
import com.freelanceplatform.android.ui.theme.FreelancePlatformTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        val tokenManager = TokenManager(this)
        var isLoggedIn by mutableStateOf(tokenManager.isLoggedIn())

        setContent {
            FreelancePlatformTheme {
                Surface(modifier = Modifier.fillMaxSize(), color = MaterialTheme.colorScheme.background) {
                    val navController = rememberNavController()
                    val navBackStackEntry by navController.currentBackStackEntryAsState()
                    val currentRoute = navBackStackEntry?.destination?.route

                    val showBottomBar = currentRoute in listOf("home", "projects", "chat", "profile")

                    androidx.compose.foundation.layout.Column {
                        if (showBottomBar) {
                            androidx.compose.foundation.layout.Column(modifier = Modifier.weight(1f)) {
                                NavHost(navController = navController, startDestination = if (isLoggedIn) "home" else "login") {
                                    composable("login") { LoginScreen(navController) { isLoggedIn = true } }
                                    composable("register") { RegisterScreen(navController) }
                                    composable("home") { HomeScreen(navController) }
                                    composable("projects") { ProjectListScreen(navController) }
                                    composable("project/{id}") { ProjectDetailScreen(navController, it.arguments?.getString("id") ?: "") }
                                    composable("create-project") { CreateProjectScreen(navController) }
                                    composable("developers") { DeveloperListScreen(navController) }
                                    composable("developer/{id}") { DeveloperDetailScreen(navController, it.arguments?.getString("id") ?: "") }
                                    composable("chat") { ChatListScreen(navController) }
                                    composable("chat/{id}") { ChatDetailScreen(navController, it.arguments?.getString("id") ?: "") }
                                    composable("profile") { ProfileScreen(navController) }
                                    composable("edit-profile") { EditProfileScreen(navController) }
                                    composable("wallet") { WalletScreen(navController) }
                                }
                            }
                            BottomNavBar(navController, currentRoute)
                        } else {
                            NavHost(navController = navController, startDestination = if (isLoggedIn) "home" else "login") {
                                composable("login") { LoginScreen(navController) { isLoggedIn = true } }
                                composable("register") { RegisterScreen(navController) }
                                composable("home") { HomeScreen(navController) }
                                composable("projects") { ProjectListScreen(navController) }
                                composable("project/{id}") { ProjectDetailScreen(navController, it.arguments?.getString("id") ?: "") }
                                composable("create-project") { CreateProjectScreen(navController) }
                                composable("developers") { DeveloperListScreen(navController) }
                                composable("developer/{id}") { DeveloperDetailScreen(navController, it.arguments?.getString("id") ?: "") }
                                composable("chat") { ChatListScreen(navController) }
                                composable("chat/{id}") { ChatDetailScreen(navController, it.arguments?.getString("id") ?: "") }
                                composable("profile") { ProfileScreen(navController) }
                                composable("edit-profile") { EditProfileScreen(navController) }
                                composable("wallet") { WalletScreen(navController) }
                            }
                        }
                    }
                }
            }
        }
    }
}
