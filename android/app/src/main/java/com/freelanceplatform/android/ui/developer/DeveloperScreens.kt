package com.freelanceplatform.android.ui.developer

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.freelanceplatform.android.data.model.DeveloperProfile
import com.freelanceplatform.android.data.api.RetrofitClient
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DeveloperListScreen(navController: NavController) {
    var developers by remember { mutableStateOf<List<DeveloperProfile>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.searchDevelopers(page = 1, pageSize = 20)
                if (res.code == 0) developers = res.data ?: emptyList()
            } catch (_: Exception) {}
            isLoading = false
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text("找开发者") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        if (isLoading) {
            Box(Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) { CircularProgressIndicator() }
        } else {
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding)) {
                items(developers) { dev ->
                    DeveloperCard(dev) { navController.navigate("developer/${dev.id}") }
                }
            }
        }
    }
}

@Composable
fun DeveloperCard(dev: DeveloperProfile, onClick: () -> Unit) {
    Card(modifier = Modifier.fillMaxWidth().padding(horizontal = 16.dp, vertical = 6.dp), shape = RoundedCornerShape(12.dp), onClick = onClick) {
        Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
            Surface(modifier = Modifier.size(48.dp), shape = CircleShape, color = MaterialTheme.colorScheme.primaryContainer) {}
            Spacer(Modifier.width(12.dp))
            Column(modifier = Modifier.weight(1f)) {
                Text(dev.title ?: "开发者", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, maxLines = 1, overflow = TextOverflow.Ellipsis)
                dev.hourlyRate?.let { Text("¥$it/小时", color = MaterialTheme.colorScheme.primary, style = MaterialTheme.typography.bodyMedium) }
                Text("⭐ ${dev.ratingAvg} · ${dev.completedProjects}个项目", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DeveloperDetailScreen(navController: NavController, developerId: String) {
    var developer by remember { mutableStateOf<DeveloperProfile?>(null) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(developerId) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.getDeveloper(developerId)
                if (res.code == 0) developer = res.data
            } catch (_: Exception) {}
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text("开发者详情") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        developer?.let { dev ->
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
                item {
                    Surface(modifier = Modifier.size(80.dp), shape = CircleShape, color = MaterialTheme.colorScheme.primaryContainer) {}
                    Spacer(Modifier.height(12.dp))
                    Text(dev.title ?: "开发者", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
                    dev.bio?.let { Text(it, style = MaterialTheme.typography.bodyLarge, modifier = Modifier.padding(top = 8.dp)) }
                    Spacer(Modifier.height(16.dp))
                    Row(horizontalArrangement = Arrangement.spacedBy(24.dp)) {
                        Column { Text("⭐ ${dev.ratingAvg}", style = MaterialTheme.typography.titleMedium); Text("评分", style = MaterialTheme.typography.bodySmall) }
                        Column { Text("${dev.completedProjects}", style = MaterialTheme.typography.titleMedium); Text("项目", style = MaterialTheme.typography.bodySmall) }
                        Column { Text("${dev.experienceYears}年", style = MaterialTheme.typography.titleMedium); Text("经验", style = MaterialTheme.typography.bodySmall) }
                        dev.hourlyRate?.let { Column { Text("¥$it", style = MaterialTheme.typography.titleMedium); Text("时薪", style = MaterialTheme.typography.bodySmall) } }
                    }
                }
            }
        }
    }
}
