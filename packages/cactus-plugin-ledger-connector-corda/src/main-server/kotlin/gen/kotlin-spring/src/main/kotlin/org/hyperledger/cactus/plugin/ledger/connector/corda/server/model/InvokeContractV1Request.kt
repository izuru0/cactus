package org.hyperledger.cactus.plugin.ledger.connector.corda.server.model

import java.util.Objects
import com.fasterxml.jackson.annotation.JsonProperty
import com.fasterxml.jackson.annotation.JsonValue
import org.hyperledger.cactus.plugin.ledger.connector.corda.server.model.FlowInvocationType
import org.hyperledger.cactus.plugin.ledger.connector.corda.server.model.JvmObject
import javax.validation.constraints.DecimalMax
import javax.validation.constraints.DecimalMin
import javax.validation.constraints.Email
import javax.validation.constraints.Max
import javax.validation.constraints.Min
import javax.validation.constraints.NotNull
import javax.validation.constraints.Pattern
import javax.validation.constraints.Size
import javax.validation.Valid
import io.swagger.v3.oas.annotations.media.Schema

/**
 * 
 * @param flowFullClassName The fully qualified name of the Corda flow to invoke
 * @param flowInvocationType 
 * @param params The list of arguments to pass in to the contract method being invoked.
 * @param timeoutMs The amount of milliseconds to wait for a transaction receipt beforegiving up and crashing.
 */
data class InvokeContractV1Request(

    @get:Size(min=1,max=1024)
    @Schema(example = "net.corda.samples.obligation.flows.IOUIssueFlow", required = true, description = "The fully qualified name of the Corda flow to invoke")
    @get:JsonProperty("flowFullClassName", required = true) val flowFullClassName: kotlin.String,

    @field:Valid
    @Schema(example = "null", required = true, description = "")
    @get:JsonProperty("flowInvocationType", required = true) val flowInvocationType: FlowInvocationType,

    @field:Valid
    @Schema(example = "null", required = true, description = "The list of arguments to pass in to the contract method being invoked.")
    @get:JsonProperty("params", required = true) val params: kotlin.collections.List<JvmObject> = arrayListOf(),

    @get:Min(0)
    @Schema(example = "null", description = "The amount of milliseconds to wait for a transaction receipt beforegiving up and crashing.")
    @get:JsonProperty("timeoutMs") val timeoutMs: kotlin.Int? = 60000
) {

}

