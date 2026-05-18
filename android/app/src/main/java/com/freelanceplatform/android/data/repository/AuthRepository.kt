package com.freelanceplatform.android.data.repository

import com.freelanceplatform.android.data.api.ApiService
import com.freelanceplatform.android.data.local.TokenManager
import com.freelanceplatform.android.data.model.*
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val apiService: ApiService,
    private val tokenManager: TokenManager
) {

    suspend fun login(email: String, password: String): Result<AuthResponse> {
        return try {
            val request = LoginRequest(email = email, password = password)
            val response = apiService.login(request)
            if (response.code == 0 && response.data != null) {
                tokenManager.saveTokens(
                    response.data.accessToken,
                    response.data.refreshToken
                )
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun loginWithPhone(phone: String, password: String): Result<AuthResponse> {
        return try {
            val request = LoginRequest(phone = phone, password = password)
            val response = apiService.login(request)
            if (response.code == 0 && response.data != null) {
                tokenManager.saveTokens(
                    response.data.accessToken,
                    response.data.refreshToken
                )
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun register(data: RegisterRequest): Result<AuthResponse> {
        return try {
            val response = apiService.register(data)
            if (response.code == 0 && response.data != null) {
                tokenManager.saveTokens(
                    response.data.accessToken,
                    response.data.refreshToken
                )
                Result.success(response.data)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun logout(): Result<Unit> {
        return try {
            apiService.logout()
            tokenManager.clearTokens()
            Result.success(Unit)
        } catch (e: Exception) {
            // Even if the API call fails, clear local tokens
            tokenManager.clearTokens()
            Result.failure(e)
        }
    }

    suspend fun refreshToken(): Result<AuthResponse> {
        return try {
            val refreshToken = tokenManager.getRefreshToken()
                ?: return Result.failure(Exception("无刷新令牌"))

            val request = RefreshTokenRequest(refreshToken)
            val response = apiService.refreshToken(request)
            if (response.code == 0 && response.data != null) {
                tokenManager.saveTokens(
                    response.data.accessToken,
                    response.data.refreshToken
                )
                Result.success(response.data)
            } else {
                tokenManager.clearTokens()
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            tokenManager.clearTokens()
            Result.failure(e)
        }
    }

    suspend fun forgotPassword(email: String): Result<Unit> {
        return try {
            val response = apiService.forgotPassword(ForgotPasswordRequest(email))
            if (response.code == 0) {
                Result.success(Unit)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun resetPassword(token: String, password: String): Result<Unit> {
        return try {
            val response = apiService.resetPassword(ResetPasswordRequest(token, password))
            if (response.code == 0) {
                Result.success(Unit)
            } else {
                Result.failure(Exception(response.message))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    fun isLoggedIn(): Boolean = tokenManager.isLoggedIn()
}
