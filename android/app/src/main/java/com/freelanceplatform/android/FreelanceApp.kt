package com.freelanceplatform.android

import android.app.Application
import dagger.hilt.android.HiltAndroidApp

@HiltAndroidApp
class FreelanceApp : Application() {

    override fun onCreate() {
        super.onCreate()
        // Hilt handles DI setup automatically via @HiltAndroidApp
    }
}
