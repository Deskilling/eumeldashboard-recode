package leif.eumeldashboard

import org.bukkit.plugin.java.JavaPlugin
import java.net.HttpURLConnection
import java.net.URI

class Eumeldashboard : JavaPlugin() {

    var serverUrl = "localhost:8080/api"
    var authToken = "leggeier"
    
    override fun onEnable() {
        // Plugin startup logic
    }

    override fun onDisable() {
        // Plugin shutdown logic
    }

    private fun sendAsync(endpoint: String, json: String) {
        server.scheduler.runTaskAsynchronously(this, Runnable {
            sendPost("$serverUrl/$endpoint", json)
        })
    }

    private fun sendPost(urlString: String, jsonBody: String) {
        var conn: HttpURLConnection? = null
        try {
            conn = URI.create(urlString).toURL().openConnection() as HttpURLConnection
            conn.apply {
                requestMethod = "POST"
                doOutput = true
                connectTimeout = 5000
                readTimeout = 5000
                setRequestProperty("Content-Type", "application/json; charset=UTF-8")
                setRequestProperty("Authorization", "Bearer $authToken")
            }
            conn.outputStream.bufferedWriter().use { it.write(jsonBody) }
            if (conn.responseCode !in 200..299) {
                logger.info("POST $urlString failed: HTTP ${conn.responseCode} - ${conn.errorStream?.bufferedReader()?.readText()}")
            }
        } catch (e: Exception) {
            logger.info("Error sending POST $urlString: ${e.message}")
        } finally {
            conn?.disconnect()
        }
    }
}
