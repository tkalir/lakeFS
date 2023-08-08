/*
 * lakeFS API
 * lakeFS HTTP API
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package io.lakefs.clients.api.model;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import org.openapitools.jackson.nullable.JsonNullable;

/**
 * location for placing an object when staging it
 */
@ApiModel(description = "location for placing an object when staging it")
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class StagingLocation {
  public static final String SERIALIZED_NAME_PHYSICAL_ADDRESS = "physical_address";
  @SerializedName(SERIALIZED_NAME_PHYSICAL_ADDRESS)
  private String physicalAddress;

  public static final String SERIALIZED_NAME_TOKEN = "token";
  @SerializedName(SERIALIZED_NAME_TOKEN)
  private String token;

  public static final String SERIALIZED_NAME_PRESIGNED_URL = "presigned_url";
  @SerializedName(SERIALIZED_NAME_PRESIGNED_URL)
  private String presignedUrl;

  public static final String SERIALIZED_NAME_PRESIGNED_URL_EXPIRY = "presigned_url_expiry";
  @SerializedName(SERIALIZED_NAME_PRESIGNED_URL_EXPIRY)
  private Long presignedUrlExpiry;


  public StagingLocation physicalAddress(String physicalAddress) {
    
    this.physicalAddress = physicalAddress;
    return this;
  }

   /**
   * Get physicalAddress
   * @return physicalAddress
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getPhysicalAddress() {
    return physicalAddress;
  }


  public void setPhysicalAddress(String physicalAddress) {
    this.physicalAddress = physicalAddress;
  }


  public StagingLocation token(String token) {
    
    this.token = token;
    return this;
  }

   /**
   * opaque staging token to use to link uploaded object
   * @return token
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "opaque staging token to use to link uploaded object")

  public String getToken() {
    return token;
  }


  public void setToken(String token) {
    this.token = token;
  }


  public StagingLocation presignedUrl(String presignedUrl) {
    
    this.presignedUrl = presignedUrl;
    return this;
  }

   /**
   * if presign&#x3D;true is passed in the request, this field will contain a presigned URL to use when uploading
   * @return presignedUrl
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "if presign=true is passed in the request, this field will contain a presigned URL to use when uploading")

  public String getPresignedUrl() {
    return presignedUrl;
  }


  public void setPresignedUrl(String presignedUrl) {
    this.presignedUrl = presignedUrl;
  }


  public StagingLocation presignedUrlExpiry(Long presignedUrlExpiry) {
    
    this.presignedUrlExpiry = presignedUrlExpiry;
    return this;
  }

   /**
   * If present and nonzero, physical_address is a presigned URL and will expire at this Unix Epoch time.  This will be shorter than the presigned URL lifetime if an authentication token is about to expire.  This field is *optional*. 
   * @return presignedUrlExpiry
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "If present and nonzero, physical_address is a presigned URL and will expire at this Unix Epoch time.  This will be shorter than the presigned URL lifetime if an authentication token is about to expire.  This field is *optional*. ")

  public Long getPresignedUrlExpiry() {
    return presignedUrlExpiry;
  }


  public void setPresignedUrlExpiry(Long presignedUrlExpiry) {
    this.presignedUrlExpiry = presignedUrlExpiry;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    StagingLocation stagingLocation = (StagingLocation) o;
    return Objects.equals(this.physicalAddress, stagingLocation.physicalAddress) &&
        Objects.equals(this.token, stagingLocation.token) &&
        Objects.equals(this.presignedUrl, stagingLocation.presignedUrl) &&
        Objects.equals(this.presignedUrlExpiry, stagingLocation.presignedUrlExpiry);
  }

  private static <T> boolean equalsNullable(JsonNullable<T> a, JsonNullable<T> b) {
    return a == b || (a != null && b != null && a.isPresent() && b.isPresent() && Objects.deepEquals(a.get(), b.get()));
  }

  @Override
  public int hashCode() {
    return Objects.hash(physicalAddress, token, presignedUrl, presignedUrlExpiry);
  }

  private static <T> int hashCodeNullable(JsonNullable<T> a) {
    if (a == null) {
      return 1;
    }
    return a.isPresent() ? Arrays.deepHashCode(new Object[]{a.get()}) : 31;
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class StagingLocation {\n");
    sb.append("    physicalAddress: ").append(toIndentedString(physicalAddress)).append("\n");
    sb.append("    token: ").append(toIndentedString(token)).append("\n");
    sb.append("    presignedUrl: ").append(toIndentedString(presignedUrl)).append("\n");
    sb.append("    presignedUrlExpiry: ").append(toIndentedString(presignedUrlExpiry)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

