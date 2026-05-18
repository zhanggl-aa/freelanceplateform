package com.freelanceplatform.android.ui.project

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.freelanceplatform.android.data.model.Project
import com.freelanceplatform.android.data.api.RetrofitClient
import com.freelanceplatform.android.ui.home.ProjectCard
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProjectListScreen(navController: NavController) {
    var projects by remember { mutableStateOf<List<Project>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var searchQuery by remember { mutableStateOf("") }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.searchProjects(page = 1, pageSize = 20, keyword = searchQuery.ifBlank { null })
                if (res.code == 0) projects = res.data ?: emptyList()
            } catch (_: Exception) {}
            isLoading = false
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text("找项目") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        if (isLoading) {
            Box(Modifier.fillMaxSize().padding(padding), contentAlignment = androidx.compose.ui.Alignment.Center) { CircularProgressIndicator() }
        } else {
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding)) {
                items(projects) { project -> ProjectCard(project) { navController.navigate("project/${project.id}") } }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProjectDetailScreen(navController: NavController, projectId: String) {
    var project by remember { mutableStateOf<Project?>(null) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(projectId) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.getProject(projectId)
                if (res.code == 0) project = res.data
            } catch (_: Exception) {}
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text(project?.title ?: "项目详情") }, navigationIcon = { Text("←", modifier = Modifier.padding(start = 16.dp).clickable { navController.popBackStack() }) }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        project?.let { p ->
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
                item {
                    Text(p.title, style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
                    Spacer(Modifier.height(8.dp))
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        p.categoryName?.let { AssistChip(onClick = {}, label = { Text(it) }) }
                        AssistChip(onClick = {}, label = { Text(p.status) })
                        AssistChip(onClick = {}, label = { Text(p.budgetType) })
                    }
                    Spacer(Modifier.height(16.dp))
                    Text("预算: ¥${p.budgetMin ?: 0} - ¥${p.budgetMax ?: 0}", style = MaterialTheme.typography.titleMedium, color = MaterialTheme.colorScheme.primary)
                    Spacer(Modifier.height(16.dp))
                    Text("项目描述", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                    Spacer(Modifier.height(8.dp))
                    Text(p.description, style = MaterialTheme.typography.bodyLarge)
                    Spacer(Modifier.height(16.dp))
                    Text("技术栈", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                    Spacer(Modifier.height(8.dp))
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        p.techStack?.forEach { tech -> AssistChip(onClick = {}, label = { Text(tech) }) }
                    }
                    Spacer(Modifier.height(24.dp))
                    Button(onClick = { /* bid dialog */ }, modifier = Modifier.fillMaxWidth()) { Text("提交投标") }
                }
            }
        }
    }
}

private fun Modifier.clickable(onClick: () -> Unit): Modifier = clickable(onClick = onClick)

@Composable
fun CreateProjectScreen(navController: NavController) {
    var title by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }
    var budgetMin by remember { mutableStateOf("") }
    var budgetMax by remember { mutableStateOf("") }

    Scaffold(topBar = {
        TopAppBar(title = { Text("发布项目") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            OutlinedTextField(value = title, onValueChange = { title = it }, label = { Text("项目标题") }, modifier = Modifier.fillMaxWidth())
            Spacer(Modifier.height(12.dp))
            OutlinedTextField(value = description, onValueChange = { description = it }, label = { Text("项目描述") }, modifier = Modifier.fillMaxWidth().height(150.dp), maxLines = 6)
            Spacer(Modifier.height(12.dp))
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                OutlinedTextField(value = budgetMin, onValueChange = { budgetMin = it }, label = { Text("最低预算") }, modifier = Modifier.weight(1f))
                OutlinedTextField(value = budgetMax, onValueChange = { budgetMax = it }, label = { Text("最高预算") }, modifier = Modifier.weight(1f))
            }
            Spacer(Modifier.height(24.dp))
            Button(onClick = { navController.popBackStack() }, modifier = Modifier.fillMaxWidth()) { Text("发布项目") }
        }
    }
}
