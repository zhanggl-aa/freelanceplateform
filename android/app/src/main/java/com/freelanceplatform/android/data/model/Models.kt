package com.freelanceplatform.android.data.model

import com.google.gson.annotations.SerializedName

// ── Domain Models ──────────────────────────────────────────────────────────

data class User(
    val id: Long,
    val email: String? = null,
    val phone: String? = null,
    val nickname: String,
    val avatarUrl: String? = null,
    val userType: String,          // "client" | "developer"
    val status: String = "active", // "active" | "suspended" | "deleted"
    val emailVerified: Boolean = false,
    val phoneVerified: Boolean = false
)

data class DeveloperProfile(
    val id: Long,
    val userId: Long,
    val title: String? = null,
    val bio: String? = null,
    val hourlyRate: Double? = null,
    val availability: String = "available", // "available" | "busy" | "unavailable"
    val experienceYears: Int = 0,
    val ratingAvg: Double = 0.0,
    val completedProjects: Int = 0,
    val skills: List<String> = emptyList(),
    val portfolio: List<PortfolioItem> = emptyList()
)

data class PortfolioItem(
    val id: Long,
    val title: String,
    val description: String? = null,
    val imageUrl: String? = null,
    val projectUrl: String? = null
)

data class ClientProfile(
    val id: Long,
    val userId: Long,
    val companyName: String? = null,
    val industry: String? = null,
    val verified: Boolean = false,
    val totalSpent: Double = 0.0
)

data class Project(
    val id: Long,
    val clientId: Long,
    val categoryId: Long,
    val title: String,
    val description: String,
    val budgetMin: Double,
    val budgetMax: Double,
    val budgetType: String, // "fixed" | "hourly"
    val deadline: String? = null,
    val techStack: List<String> = emptyList(),
    val status: String, // "draft" | "open" | "in_progress" | "completed" | "cancelled" | "disputed"
    val viewCount: Int = 0,
    val bidCount: Int = 0,
    val categoryName: String? = null,
    val clientName: String? = null,
    val createdAt: String? = null,
    val updatedAt: String? = null
)

data class Bid(
    val id: Long,
    val projectId: Long,
    val developerId: Long,
    val coverLetter: String,
    val estimatedDays: Int,
    val proposedBudget: Double,
    val status: String, // "pending" | "shortlisted" | "accepted" | "rejected" | "countered" | "withdrawn"
    val developerName: String? = null,
    val developerAvatarUrl: String? = null,
    val createdAt: String? = null,
    val updatedAt: String? = null
)

data class Contract(
    val id: Long,
    val projectId: Long,
    val clientId: Long,
    val developerId: Long,
    val totalAmount: Double,
    val platformFee: Double,
    val developerPayout: Double,
    val paidAmount: Double = 0.0,
    val releasedAmount: Double = 0.0,
    val status: String, // "pending" | "active" | "completed" | "cancelled" | "disputed"
    val createdAt: String? = null
)

data class Milestone(
    val id: Long,
    val projectId: Long,
    val title: String,
    val amount: Double,
    val deadline: String? = null,
    val status: String, // "pending" | "submitted" | "approved" | "rejected" | "disputed" | "paid"
    val sortOrder: Int = 0
)

data class Payment(
    val id: Long? = null,
    val contractId: Long,
    val amount: Double,
    val platformFee: Double,
    val netAmount: Double,
    val paymentMethod: String? = null,
    val status: String, // "pending" | "completed" | "failed" | "refunded"
    val createdAt: String? = null
)

data class Conversation(
    val id: Long,
    val type: String = "direct", // "direct" | "project"
    val projectId: Long? = null,
    val lastMessageAt: String? = null,
    val participants: List<Participant> = emptyList(),
    val lastMessage: ChatMessage? = null,
    val unreadCount: Int = 0
)

data class Participant(
    val userId: Long,
    val nickname: String,
    val avatarUrl: String? = null
)

data class ChatMessage(
    val id: Long,
    val conversationId: Long,
    val senderId: Long,
    val content: String,
    val messageType: String = "text", // "text" | "image" | "file" | "system"
    val fileUrl: String? = null,
    val createdAt: String,
    val senderName: String? = null,
    val senderAvatarUrl: String? = null
)

data class Review(
    val id: Long,
    val projectId: Long,
    val reviewerId: Long,
    val revieweeId: Long,
    val qualityRating: Int,
    val communicationRating: Int,
    val timelinessRating: Int,
    val overallRating: Int,
    val comment: String? = null,
    val createdAt: String? = null
)

data class Notification(
    val id: Long,
    val userId: Long,
    val type: String,
    val title: String,
    val content: String,
    val isRead: Boolean = false,
    val createdAt: String
)

data class Category(
    val id: Long,
    val name: String,
    val icon: String? = null,
    val description: String? = null,
    val parentId: Long? = null,
    val projectCount: Int = 0
)

data class Wallet(
    val balance: Double = 0.0,
    val frozenAmount: Double = 0.0,
    val totalDeposited: Double = 0.0,
    val totalWithdrawn: Double = 0.0
)

data class WalletTransaction(
    val id: Long,
    val type: String, // "deposit" | "withdraw" | "payment" | "release" | "refund" | "fee"
    val amount: Double,
    val description: String? = null,
    val status: String = "completed",
    val createdAt: String
)

data class Skill(
    val id: Long? = null,
    val name: String
)

data class FileUploadResult(
    val id: Long,
    val url: String,
    val fileName: String,
    val fileSize: Long,
    val mimeType: String
)

// ── API Wrapper ────────────────────────────────────────────────────────────

data class ApiResponse<T>(
    val code: Int = 0,
    val message: String = "",
    val data: T? = null,
    val meta: Meta? = null
)

data class Meta(
    val page: Int = 1,
    val pageSize: Int = 20,
    val total: Int = 0,
    val totalPages: Int = 0
)

// ── Auth Models ────────────────────────────────────────────────────────────

data class AuthResponse(
    val accessToken: String,
    val refreshToken: String,
    val user: User? = null
)

data class LoginRequest(
    val email: String? = null,
    val phone: String? = null,
    val password: String
)

data class RegisterRequest(
    val email: String? = null,
    val phone: String? = null,
    val password: String,
    val nickname: String,
    val userType: String // "client" | "developer"
)

// ── Request Models ─────────────────────────────────────────────────────────

data class CreateProjectRequest(
    val title: String,
    val categoryId: Long,
    val description: String,
    val budgetMin: Double,
    val budgetMax: Double,
    val budgetType: String,
    val deadline: String? = null,
    val techStack: List<String> = emptyList()
)

data class UpdateProjectRequest(
    val title: String? = null,
    val categoryId: Long? = null,
    val description: String? = null,
    val budgetMin: Double? = null,
    val budgetMax: Double? = null,
    val budgetType: String? = null,
    val deadline: String? = null,
    val techStack: List<String>? = null
)

data class CreateBidRequest(
    val projectId: Long,
    val coverLetter: String,
    val estimatedDays: Int,
    val proposedBudget: Double,
    val budgetType: String = "fixed"
)

data class UpdateBidRequest(
    val coverLetter: String? = null,
    val estimatedDays: Int? = null,
    val proposedBudget: Double? = null
)

data class CounterOfferRequest(
    val proposedBudget: Double,
    val estimatedDays: Int,
    val message: String? = null
)

data class CreateContractRequest(
    val projectId: Long,
    val bidId: Long,
    val milestones: List<CreateMilestoneRequest>? = null
)

data class CreateMilestoneRequest(
    val title: String,
    val amount: Double,
    val deadline: String? = null,
    val sortOrder: Int = 0
)

data class UpdateDeveloperProfileRequest(
    val title: String? = null,
    val bio: String? = null,
    val hourlyRate: Double? = null,
    val availability: String? = null,
    val experienceYears: Int? = null,
    val skills: List<String>? = null
)

data class UpdateClientProfileRequest(
    val companyName: String? = null,
    val industry: String? = null
)

data class CreateReviewRequest(
    val projectId: Long,
    val revieweeId: Long,
    val qualityRating: Int,
    val communicationRating: Int,
    val timelinessRating: Int,
    val overallRating: Int,
    val comment: String? = null
)

data class DepositRequest(
    val amount: Double,
    val paymentMethod: String
)

data class WithdrawRequest(
    val amount: Double,
    val paymentMethod: String
)

data class ReleasePaymentRequest(
    val milestoneId: Long? = null,
    val amount: Double
)

data class ForgotPasswordRequest(
    val email: String
)

data class ResetPasswordRequest(
    val token: String,
    val password: String
)

data class RefreshTokenRequest(
    val refreshToken: String
)

data class SendMessageRequest(
    val content: String,
    val messageType: String = "text",
    val fileUrl: String? = null
)
