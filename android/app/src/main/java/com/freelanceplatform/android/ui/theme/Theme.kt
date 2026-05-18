package com.freelanceplatform.android.ui.theme

import android.os.Build
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

val Primary = Color(0xFF409EFF)
val PrimaryDark = Color(0xFF337ECC)
val PrimaryLight = Color(0xFF79BBFF)
val Success = Color(0xFF67C23A)
val Warning = Color(0xFFE6A23C)
val Danger = Color(0xFFF56C6C)
val Info = Color(0xFF909399)
val TextPrimary = Color(0xFF303133)
val TextRegular = Color(0xFF606266)
val TextSecondary = Color(0xFF909399)
val BgPage = Color(0xFFF5F7FA)
val BorderColor = Color(0xFFDCDFE6)

private val LightColorScheme = lightColorScheme(
    primary = Primary,
    onPrimary = Color.White,
    primaryContainer = PrimaryLight,
    secondary = Success,
    tertiary = Warning,
    error = Danger,
    background = BgPage,
    surface = Color.White,
    onBackground = TextPrimary,
    onSurface = TextPrimary,
    outline = BorderColor,
)

private val DarkColorScheme = darkColorScheme(
    primary = PrimaryLight,
    onPrimary = Color.White,
    secondary = Success,
    tertiary = Warning,
    error = Danger,
)

@Composable
fun FreelancePlatformTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    dynamicColor: Boolean = true,
    content: @Composable () -> Unit,
) {
    val colorScheme = when {
        dynamicColor && Build.VERSION.SDK_INT >= Build.VERSION_CODES.S -> {
            val context = LocalContext.current
            if (darkTheme) dynamicDarkColorScheme(context) else dynamicLightColorScheme(context)
        }
        darkTheme -> DarkColorScheme
        else -> LightColorScheme
    }

    val typography = Typography(
        headlineLarge = TextStyle(fontSize = 28.sp, fontWeight = FontWeight.Bold, color = TextPrimary),
        headlineMedium = TextStyle(fontSize = 24.sp, fontWeight = FontWeight.SemiBold, color = TextPrimary),
        titleLarge = TextStyle(fontSize = 20.sp, fontWeight = FontWeight.Medium, color = TextPrimary),
        titleMedium = TextStyle(fontSize = 16.sp, fontWeight = FontWeight.Medium, color = TextPrimary),
        bodyLarge = TextStyle(fontSize = 16.sp, color = TextRegular),
        bodyMedium = TextStyle(fontSize = 14.sp, color = TextRegular),
        bodySmall = TextStyle(fontSize = 12.sp, color = TextSecondary),
        labelLarge = TextStyle(fontSize = 14.sp, fontWeight = FontWeight.Medium, color = Primary),
        labelSmall = TextStyle(fontSize = 10.sp, color = TextSecondary),
    )

    MaterialTheme(
        colorScheme = colorScheme,
        typography = typography,
        content = content,
    )
}
