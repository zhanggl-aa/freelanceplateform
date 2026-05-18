package com.freelanceplatform.android.ui.home

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.freelanceplatform.android.data.model.Project
import com.freelanceplatform.android.data.api.RetrofitClient
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(navController: NavController) {
    var searchQuery by remember { mutableStateOf("") }
    var featuredProjects by remember { mutableStateOf<List<Project>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.searchProjects(page = 1, pageSize = 10)
                if (res.code == 0) featuredProjects = res.data ?: emptyList()
            } catch (_: Exception) {}
            isLoading = false
        }
    }

    val categories = listOf("网站开发", "移动应用", "小程序", "前端", "后端", "AI/ML", "UI设计", "更多")

    Scaffold(topBar = {
        TopAppBar(title = { Text("接单平台", fontWeight = FontWeight.Bold) }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(padding)) {
            item {
                OutlinedTextField(value = searchQuery, onValueChange = { searchQuery = it }, modifier = Modifier.fillMaxWidth().padding(16.dp), placeholder = { Text("搜索项目或技能") }, leadingIcon = { Icon(Icons.Default.Search, null) }, shape = RoundedCornerShape(24.dp), singleLine = true)
            }
            item {
                Text("项目分类", style = MaterialTheme.typography.titleMedium, modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp))
                LazyRow(contentPadding = PaddingValues(horizontal = 16.dp), horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                    items(categories) { cat ->
                        AssistChip(onClick = { navController.navigate("projects") }, label = { Text(cat) })
                    }
                }
            }
            item {
                Text("精选项目", style = MaterialTheme.typography.titleMedium, modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp))
            }
            if (isLoading) {
                item { Box(Modifier.fillMaxWidth().padding(32.dp), contentAlignment = Alignment.Center) { CircularProgressIndicator() } }
            } else {
                items(featuredProjects) { project ->
                    ProjectCard(project) { navController.navigate("project/${project.id}") }
                }
            }
        }
    }
}

@Composable
fun ProjectCard(project: Project, onClick: () -> Unit) {
    Card(modifier = Modifier.fillMaxWidth().padding(horizontal = 16.dp, vertical = 6.dp).clickable(onClick = onClick), shape = RoundedCornerShape(12.dp)) {
        Column(modifier = Modifier.padding(16.dp)) {
            Text(project.title, style = MaterialTheme.typography.titleMedium, maxLines = 1, overflow = TextOverflow.Ellipsis)
            Spacer(Modifier.height(8.dp))
            Row(horizontalArrangement = Arrangement.SpaceBetween, modifier = Modifier.fillMaxWidth()) {
                Text("¥${project.budgetMin ?: 0} - ¥${project.budgetMax ?: 0}", color = MaterialTheme.colorScheme.primary, style = MaterialTheme.typography.labelLarge)
                Text("${project.bidCount ?: 0}个投标", color = MaterialTheme.colorScheme.onSurfaceVariant, style = MaterialTheme.typography.bodySmall)
            }
            project.categoryName?.let {
                Spacer(Modifier.height(4.dp))
                AssistChip(onClick = {}, label = { Text(it, style = MaterialTheme.typography.bodySmall) })
            }
        }
    }
}
