package com.freelanceplatform.android.data.repository

import com.freelanceplatform.android.data.api.ApiService
import com.freelanceplatform.android.data.model.*
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ProjectRepository @Inject constructor(
    private val apiService: ApiService
) {

    suspend fun getProjects(
        page: Int = 1,
        pageSize: Int = 20,
        categoryId: Long? = null,
        status: String? = null,
        budgetType: String? = null,
        minBudget: Double? = null,
        maxBudget: Double? = null,
        search: String? = null,
        sort: String? = null,
        techStack: String? = null
    ): Result<Pair<List<Project>, Meta?>> {
        return try {
            val response = apiService.getProjects(
                page = page,
                pageSize = pageSize,
                categoryId = categoryId,
                status = status,
                budgetType = budgetType,
                minBudget = minBudget,
                maxBudget = maxBudget,
                search = search,
                sort = sort,
                techStack = techStack
            )
            if (response.code == 0 && response.data != null) {
                Result.success(response.data to response.meta)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProject(id: Long): Result<Project> {
        return try {
            val response = apiService.getProject(id)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun createProject(request: CreateProjectRequest): Result<Project> {
        return try {
            val response = apiService.createProject(request)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun updateProject(id: Long, request: UpdateProjectRequest): Result<Project> {
        return try {
            val response = apiService.updateProject(id, request)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun deleteProject(id: Long): Result<Unit> {
        return try {
            val response = apiService.deleteProject(id)
            if (response.code == 0) {
                Result.success(Unit)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun publishProject(id: Long): Result<Project> {
        return try {
            val response = apiService.publishProject(id)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun closeProject(id: Long): Result<Project> {
        return try {
            val response = apiService.closeProject(id)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun bookmarkProject(id: Long): Result<Unit> {
        return try {
            val response = apiService.bookmarkProject(id)
            if (response.code == 0) {
                Result.success(Unit)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun unbookmarkProject(id: Long): Result<Unit> {
        return try {
            val response = apiService.unbookmarkProject(id)
            if (response.code == 0) {
                Result.success(Unit)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getBookmarkedProjects(page: Int = 1, pageSize: Int = 20): Result<Pair<List<Project>, Meta?>> {
        return try {
            val response = apiService.getBookmarkedProjects(page, pageSize)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data to response.meta)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getMyBids(page: Int = 1, pageSize: Int = 20, status: String? = null): Result<Pair<List<Bid>, Meta?>> {
        return try {
            val response = apiService.getMyBids(page, pageSize, status)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data to response.meta)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun createBid(request: CreateBidRequest): Result<Bid> {
        return try {
            val response = apiService.createBid(request.projectId, request)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getProjectBids(projectId: Long, page: Int = 1, pageSize: Int = 20): Result<Pair<List<Bid>, Meta?>> {
        return try {
            val response = apiService.getProjectBids(projectId, page, pageSize)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data to response.meta)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun acceptBid(bidId: Long): Result<Bid> {
        return try {
            val response = apiService.acceptBid(bidId)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun rejectBid(bidId: Long): Result<Bid> {
        return try {
            val response = apiService.rejectBid(bidId)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getCategories(): Result<List<Category>> {
        return try {
            val response = apiService.getCategories()
            if (response.code == 0 && response.data != null) {
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getMyContracts(page: Int = 1, pageSize: Int = 20, status: String? = null): Result<Pair<List<Contract>, Meta?>> {
        return try {
            val response = apiService.getMyContracts(page, pageSize, status)
            if (response.code == 0 && response.data != null) {
                Result.success(response.data to response.meta)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
