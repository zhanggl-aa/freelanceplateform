package com.freelanceplatform.android.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.freelanceplatform.android.data.repository.AuthRepository
import com.freelanceplatform.android.data.model.RegisterRequest
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class AuthUiState(
    val isLoading: Boolean = false,
    val error: String? = null,
    val loginSuccess: Boolean = false,
    val registerSuccess: Boolean = false
)

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    fun login(email: String, password: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            val result = authRepository.login(email, password)
            _uiState.value = if (result.isSuccess) {
                _uiState.value.copy(isLoading = false, loginSuccess = true)
            } else {
                _uiState.value.copy(
                    isLoading = false,
                    error = result.exceptionOrNull()?.message ?: "登录失败"
                )
            }
        }
    }

    fun loginWithPhone(phone: String, password: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            val result = authRepository.loginWithPhone(phone, password)
            _uiState.value = if (result.isSuccess) {
                _uiState.value.copy(isLoading = false, loginSuccess = true)
            } else {
                _uiState.value.copy(
                    isLoading = false,
                    error = result.exceptionOrNull()?.message ?: "登录失败"
                )
            }
        }
    }

    fun register(email: String, phone: String, password: String, nickname: String, userType: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            val request = RegisterRequest(
                email = email.ifBlank { null },
                phone = phone.ifBlank { null },
                password = password,
                nickname = nickname,
                userType = userType
            )
            val result = authRepository.register(request)
            _uiState.value = if (result.isSuccess) {
                _uiState.value.copy(isLoading = false, registerSuccess = true, loginSuccess = true)
            } else {
                _uiState.value.copy(
                    isLoading = false,
                    error = result.exceptionOrNull()?.message ?: "注册失败"
                )
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }

    fun resetState() {
        _uiState.value = AuthUiState()
    }
}
