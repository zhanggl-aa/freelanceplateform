package com.freelanceplatform.android.data.api

import com.freelanceplatform.android.data.model.*
import okhttp3.MultipartBody
import retrofit2.http.*

interface ApiService {

    // ── Auth ────────────────────────────────────────────────────────────────

    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): ApiResponse<AuthResponse>

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): ApiResponse<AuthResponse>

    @POST("auth/logout")
    suspend fun logout(): ApiResponse<Unit>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body request: RefreshTokenRequest): ApiResponse<AuthResponse>

    @POST("auth/forgot-password")
    suspend fun forgotPassword(@Body request: ForgotPasswordRequest): ApiResponse<Unit>

    @POST("auth/reset-password")
    suspend fun resetPassword(@Body request: ResetPasswordRequest): ApiResponse<Unit>

    // ── Users ───────────────────────────────────────────────────────────────

    @GET("users/me")
    suspend fun getMyProfile(): ApiResponse<User>

    @PUT("users/me")
    suspend fun updateMyProfile(@Body user: User): ApiResponse<User>

    @DELETE("users/me")
    suspend fun deleteMyAccount(): ApiResponse<Unit>

    @GET("users/{id}")
    suspend fun getUser(@Path("id") id: Long): ApiResponse<User>

    // ── Developer Profile ───────────────────────────────────────────────────

    @POST("developers/profile")
    suspend fun createDeveloperProfile(@Body request: UpdateDeveloperProfileRequest): ApiResponse<DeveloperProfile>

    @GET("developers/profile")
    suspend fun getMyDeveloperProfile(): ApiResponse<DeveloperProfile>

    @PUT("developers/profile")
    suspend fun updateDeveloperProfile(@Body request: UpdateDeveloperProfileRequest): ApiResponse<DeveloperProfile>

    @GET("developers")
    suspend fun getDevelopers(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20,
        @Query("skills") skills: String? = null,
        @Query("minRate") minRate: Double? = null,
        @Query("maxRate") maxRate: Double? = null,
        @Query("availability") availability: String? = null,
        @Query("search") search: String? = null
    ): ApiResponse<List<DeveloperProfile>>

    @GET("developers/{id}")
    suspend fun getDeveloper(@Path("id") id: Long): ApiResponse<DeveloperProfile>

    @DELETE("developers/profile")
    suspend fun deleteDeveloperProfile(): ApiResponse<Unit>

    // ── Developer Skills ────────────────────────────────────────────────────

    @POST("developers/profile/skills")
    suspend fun addSkill(@Body skill: Skill): ApiResponse<List<Skill>>

    @DELETE("developers/profile/skills/{skillId}")
    suspend fun removeSkill(@Path("skillId") skillId: Long): ApiResponse<List<Skill>>

    // ── Developer Portfolio ─────────────────────────────────────────────────

    @POST("developers/profile/portfolio")
    suspend fun addPortfolioItem(@Body item: PortfolioItem): ApiResponse<PortfolioItem>

    @PUT("developers/profile/portfolio/{id}")
    suspend fun updatePortfolioItem(@Path("id") id: Long, @Body item: PortfolioItem): ApiResponse<PortfolioItem>

    @DELETE("developers/profile/portfolio/{id}")
    suspend fun deletePortfolioItem(@Path("id") id: Long): ApiResponse<Unit>

    // ── Client Profile ──────────────────────────────────────────────────────

    @POST("clients/profile")
    suspend fun createClientProfile(@Body request: UpdateClientProfileRequest): ApiResponse<ClientProfile>

    @GET("clients/profile")
    suspend fun getMyClientProfile(): ApiResponse<ClientProfile>

    @PUT("clients/profile")
    suspend fun updateClientProfile(@Body request: UpdateClientProfileRequest): ApiResponse<ClientProfile>

    @DELETE("clients/profile")
    suspend fun deleteClientProfile(): ApiResponse<Unit>

    // ── Categories ──────────────────────────────────────────────────────────

    @GET("categories")
    suspend fun getCategories(): ApiResponse<List<Category>>

    // ── Projects ────────────────────────────────────────────────────────────

    @POST("projects")
    suspend fun createProject(@Body request: CreateProjectRequest): ApiResponse<Project>

    @GET("projects")
    suspend fun getProjects(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20,
        @Query("categoryId") categoryId: Long? = null,
        @Query("status") status: String? = null,
        @Query("budgetType") budgetType: String? = null,
        @Query("minBudget") minBudget: Double? = null,
        @Query("maxBudget") maxBudget: Double? = null,
        @Query("search") search: String? = null,
        @Query("sort") sort: String? = null, // "newest" | "budget_high" | "budget_low" | "bids"
        @Query("techStack") techStack: String? = null
    ): ApiResponse<List<Project>>

    @GET("projects/{id}")
    suspend fun getProject(@Path("id") id: Long): ApiResponse<Project>

    @PUT("projects/{id}")
    suspend fun updateProject(
        @Path("id") id: Long,
        @Body request: UpdateProjectRequest
    ): ApiResponse<Project>

    @DELETE("projects/{id}")
    suspend fun deleteProject(@Path("id") id: Long): ApiResponse<Unit>

    @POST("projects/{id}/publish")
    suspend fun publishProject(@Path("id") id: Long): ApiResponse<Project>

    @POST("projects/{id}/close")
    suspend fun closeProject(@Path("id") id: Long): ApiResponse<Project>

    // ── Project Bookmarks ───────────────────────────────────────────────────

    @POST("projects/{id}/bookmark")
    suspend fun bookmarkProject(@Path("id") id: Long): ApiResponse<Unit>

    @DELETE("projects/{id}/bookmark")
    suspend fun unbookmarkProject(@Path("id") id: Long): ApiResponse<Unit>

    @GET("users/me/bookmarks")
    suspend fun getBookmarkedProjects(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20
    ): ApiResponse<List<Project>>

    // ── Bids ────────────────────────────────────────────────────────────────

    @POST("projects/{projectId}/bids")
    suspend fun createBid(
        @Path("projectId") projectId: Long,
        @Body request: CreateBidRequest
    ): ApiResponse<Bid>

    @GET("projects/{projectId}/bids")
    suspend fun getProjectBids(
        @Path("projectId") projectId: Long,
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20
    ): ApiResponse<List<Bid>>

    @GET("bids/{id}")
    suspend fun getBid(@Path("id") id: Long): ApiResponse<Bid>

    @PUT("bids/{id}")
    suspend fun updateBid(
        @Path("id") id: Long,
        @Body request: UpdateBidRequest
    ): ApiResponse<Bid>

    @DELETE("bids/{id}")
    suspend fun withdrawBid(@Path("id") id: Long): ApiResponse<Unit>

    @POST("bids/{id}/accept")
    suspend fun acceptBid(@Path("id") id: Long): ApiResponse<Bid>

    @POST("bids/{id}/reject")
    suspend fun rejectBid(@Path("id") id: Long): ApiResponse<Bid>

    @POST("bids/{id}/shortlist")
    suspend fun shortlistBid(@Path("id") id: Long): ApiResponse<Bid>

    @POST("bids/{id}/counter-offer")
    suspend fun counterOfferBid(
        @Path("id") id: Long,
        @Body request: CounterOfferRequest
    ): ApiResponse<Bid>

    @GET("users/me/bids")
    suspend fun getMyBids(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20,
        @Query("status") status: String? = null
    ): ApiResponse<List<Bid>>

    // ── Contracts ───────────────────────────────────────────────────────────

    @GET("contracts/my")
    suspend fun getMyContracts(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20,
        @Query("status") status: String? = null
    ): ApiResponse<List<Contract>>

    @GET("contracts/{id}")
    suspend fun getContract(@Path("id") id: Long): ApiResponse<Contract>

    @POST("contracts/{id}/start")
    suspend fun startContract(@Path("id") id: Long): ApiResponse<Contract>

    @POST("contracts/{id}/cancel")
    suspend fun cancelContract(@Path("id") id: Long): ApiResponse<Contract>

    @POST("contracts/{id}/dispute")
    suspend fun disputeContract(@Path("id") id: Long): ApiResponse<Contract>

    // ── Milestones ──────────────────────────────────────────────────────────

    @POST("contracts/{contractId}/milestones")
    suspend fun createMilestone(
        @Path("contractId") contractId: Long,
        @Body request: CreateMilestoneRequest
    ): ApiResponse<Milestone>

    @GET("contracts/{contractId}/milestones")
    suspend fun getMilestones(@Path("contractId") contractId: Long): ApiResponse<List<Milestone>>

    @PUT("milestones/{id}")
    suspend fun updateMilestone(
        @Path("id") id: Long,
        @Body request: CreateMilestoneRequest
    ): ApiResponse<Milestone>

    @DELETE("milestones/{id}")
    suspend fun deleteMilestone(@Path("id") id: Long): ApiResponse<Unit>

    @POST("milestones/{id}/submit")
    suspend fun submitMilestone(@Path("id") id: Long): ApiResponse<Milestone>

    @POST("milestones/{id}/approve")
    suspend fun approveMilestone(@Path("id") id: Long): ApiResponse<Milestone>

    @POST("milestones/{id}/reject")
    suspend fun rejectMilestone(@Path("id") id: Long): ApiResponse<Milestone>

    @POST("milestones/{id}/dispute")
    suspend fun disputeMilestone(@Path("id") id: Long): ApiResponse<Milestone>

    // ── Payments ────────────────────────────────────────────────────────────

    @POST("payments/deposit")
    suspend fun deposit(@Body request: DepositRequest): ApiResponse<Payment>

    @POST("payments/release")
    suspend fun releasePayment(@Body request: ReleasePaymentRequest): ApiResponse<Payment>

    @POST("payments/refund")
    suspend fun refundPayment(
        @Query("contractId") contractId: Long,
        @Query("amount") amount: Double
    ): ApiResponse<Payment>

    @GET("wallet")
    suspend fun getWallet(): ApiResponse<Wallet>

    @GET("wallet/transactions")
    suspend fun getWalletTransactions(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20
    ): ApiResponse<List<WalletTransaction>>

    @POST("wallet/withdraw")
    suspend fun withdraw(@Body request: WithdrawRequest): ApiResponse<Payment>

    // ── Chat / Conversations ────────────────────────────────────────────────

    @GET("conversations")
    suspend fun getConversations(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20
    ): ApiResponse<List<Conversation>>

    @GET("conversations/{id}")
    suspend fun getConversation(@Path("id") id: Long): ApiResponse<Conversation>

    @POST("conversations")
    suspend fun createConversation(
        @Query("participantId") participantId: Long,
        @Query("projectId") projectId: Long? = null
    ): ApiResponse<Conversation>

    @DELETE("conversations/{id}")
    suspend fun deleteConversation(@Path("id") id: Long): ApiResponse<Unit>

    @GET("conversations/{id}/messages")
    suspend fun getMessages(
        @Path("id") conversationId: Long,
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 50,
        @Query("before") before: String? = null
    ): ApiResponse<List<ChatMessage>>

    @POST("conversations/{id}/messages")
    suspend fun sendMessage(
        @Path("id") conversationId: Long,
        @Body request: SendMessageRequest
    ): ApiResponse<ChatMessage>

    // ── Reviews ─────────────────────────────────────────────────────────────

    @POST("reviews")
    suspend fun createReview(@Body request: CreateReviewRequest): ApiResponse<Review>

    @GET("reviews")
    suspend fun getReviews(
        @Query("userId") userId: Long,
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20
    ): ApiResponse<List<Review>>

    @GET("reviews/{id}")
    suspend fun getReview(@Path("id") id: Long): ApiResponse<Review>

    // ── Notifications ───────────────────────────────────────────────────────

    @GET("notifications")
    suspend fun getNotifications(
        @Query("page") page: Int = 1,
        @Query("pageSize") pageSize: Int = 20
    ): ApiResponse<List<Notification>>

    @GET("notifications/unread-count")
    suspend fun getUnreadNotificationCount(): ApiResponse<Int>

    @PUT("notifications/{id}/read")
    suspend fun markNotificationRead(@Path("id") id: Long): ApiResponse<Unit>

    @PUT("notifications/read-all")
    suspend fun markAllNotificationsRead(): ApiResponse<Unit>

    // ── File Upload ─────────────────────────────────────────────────────────

    @Multipart
    @POST("files/upload")
    suspend fun uploadFile(
        @Part file: MultipartBody.Part,
        @Query("type") type: String = "general"
    ): ApiResponse<FileUploadResult>

    @GET("files/{id}")
    suspend fun getFile(@Path("id") id: Long): ApiResponse<FileUploadResult>

    @DELETE("files/{id}")
    suspend fun deleteFile(@Path("id") id: Long): ApiResponse<Unit>
}
