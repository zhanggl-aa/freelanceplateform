package com.freelanceplatform.android.data.api

import com.freelanceplatform.android.data.local.TokenManager
import com.freelanceplatform.android.data.model.ApiResponse
import com.freelanceplatform.android.data.model.AuthResponse
import com.freelanceplatform.android.data.model.RefreshTokenRequest
import com.google.gson.Gson
import com.google.gson.reflect.TypeToken
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import okhttp3.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TokenInterceptor @Inject constructor(
    private val tokenManager: TokenManager
) : Interceptor {

    private val gson = Gson()

    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()

        // Skip token for auth endpoints
        val path = originalRequest.url.encodedPath
        if (path.contains("auth/login") ||
            path.contains("auth/register") ||
            path.contains("auth/forgot-password") ||
            path.contains("auth/reset-password")
        ) {
            return chain.proceed(originalRequest)
        }

        val accessToken = tokenManager.getAccessToken()
        val request = if (!accessToken.isNullOrEmpty()) {
            originalRequest.newBuilder()
                .header("Authorization", "Bearer $accessToken")
                .build()
        } else {
            originalRequest
        }

        val response = chain.proceed(request)

        // On 401, try to refresh the token and retry
        if (response.code == 401) {
            response.close()

            val refreshToken = tokenManager.getRefreshToken()
            if (refreshToken.isNullOrEmpty()) {
                tokenManager.clearTokens()
                return response
            }

            val newTokens = refreshTokens(chain, refreshToken)
            if (newTokens != null) {
                tokenManager.saveTokens(newTokens.accessToken, newTokens.refreshToken)

                // Retry the original request with new token
                val retriedRequest = originalRequest.newBuilder()
                    .header("Authorization", "Bearer ${newTokens.accessToken}")
                    .build()
                return chain.proceed(retriedRequest)
            } else {
                // Refresh failed — clear tokens and force re-login
                tokenManager.clearTokens()
                return response
            }
        }

        return response
    }

    private fun refreshTokens(chain: Interceptor.Chain, refreshToken: String): AuthResponse? {
        return try {
            val baseUrl = chain.request().url.run {
                "${scheme}://${host}:${port}/api/v1/"
            }
            val json = gson.toJson(RefreshTokenRequest(refreshToken))
            val body = json.toRequestBody("application/json; charset=utf-8".toMediaType())

            val refreshRequest = Request.Builder()
                .url("${baseUrl}auth/refresh")
                .post(body)
                .build()

            // Use a separate client for refresh to avoid interceptor loop
            val refreshClient = OkHttpClient.Builder().build()
            val refreshResponse = refreshClient.newCall(refreshRequest).execute()

            if (refreshResponse.isSuccessful) {
                val responseBody = refreshResponse.body?.string() ?: return null
                val type = object : TypeToken<ApiResponse<AuthResponse>>() {}.type
                val apiResponse: ApiResponse<AuthResponse> = gson.fromJson(responseBody, type)
                apiResponse.data
            } else {
                null
            }
        } catch (e: Exception) {
            null
        }
    }
}
