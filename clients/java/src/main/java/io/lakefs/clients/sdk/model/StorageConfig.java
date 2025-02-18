/*
 * lakeFS API
 * lakeFS HTTP API
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package io.lakefs.clients.sdk.model;

import java.util.Objects;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;
import java.util.Arrays;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonArray;
import com.google.gson.JsonDeserializationContext;
import com.google.gson.JsonDeserializer;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParseException;
import com.google.gson.TypeAdapterFactory;
import com.google.gson.reflect.TypeToken;
import com.google.gson.TypeAdapter;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;

import java.lang.reflect.Type;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

import io.lakefs.clients.sdk.JSON;

/**
 * StorageConfig
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class StorageConfig {
  public static final String SERIALIZED_NAME_BLOCKSTORE_TYPE = "blockstore_type";
  @SerializedName(SERIALIZED_NAME_BLOCKSTORE_TYPE)
  private String blockstoreType;

  public static final String SERIALIZED_NAME_BLOCKSTORE_NAMESPACE_EXAMPLE = "blockstore_namespace_example";
  @SerializedName(SERIALIZED_NAME_BLOCKSTORE_NAMESPACE_EXAMPLE)
  private String blockstoreNamespaceExample;

  public static final String SERIALIZED_NAME_BLOCKSTORE_NAMESPACE_VALIDITY_REGEX = "blockstore_namespace_ValidityRegex";
  @SerializedName(SERIALIZED_NAME_BLOCKSTORE_NAMESPACE_VALIDITY_REGEX)
  private String blockstoreNamespaceValidityRegex;

  public static final String SERIALIZED_NAME_DEFAULT_NAMESPACE_PREFIX = "default_namespace_prefix";
  @SerializedName(SERIALIZED_NAME_DEFAULT_NAMESPACE_PREFIX)
  private String defaultNamespacePrefix;

  public static final String SERIALIZED_NAME_PRE_SIGN_SUPPORT = "pre_sign_support";
  @SerializedName(SERIALIZED_NAME_PRE_SIGN_SUPPORT)
  private Boolean preSignSupport;

  public static final String SERIALIZED_NAME_PRE_SIGN_SUPPORT_UI = "pre_sign_support_ui";
  @SerializedName(SERIALIZED_NAME_PRE_SIGN_SUPPORT_UI)
  private Boolean preSignSupportUi;

  public static final String SERIALIZED_NAME_IMPORT_SUPPORT = "import_support";
  @SerializedName(SERIALIZED_NAME_IMPORT_SUPPORT)
  private Boolean importSupport;

  public static final String SERIALIZED_NAME_IMPORT_VALIDITY_REGEX = "import_validity_regex";
  @SerializedName(SERIALIZED_NAME_IMPORT_VALIDITY_REGEX)
  private String importValidityRegex;

  public static final String SERIALIZED_NAME_PRE_SIGN_MULTIPART_UPLOAD = "pre_sign_multipart_upload";
  @SerializedName(SERIALIZED_NAME_PRE_SIGN_MULTIPART_UPLOAD)
  private Boolean preSignMultipartUpload;

  public static final String SERIALIZED_NAME_BLOCKSTORE_ID = "blockstore_id";
  @SerializedName(SERIALIZED_NAME_BLOCKSTORE_ID)
  private String blockstoreId;

  public static final String SERIALIZED_NAME_BLOCKSTORE_DESCRIPTION = "blockstore_description";
  @SerializedName(SERIALIZED_NAME_BLOCKSTORE_DESCRIPTION)
  private String blockstoreDescription;

  public StorageConfig() {
  }

  public StorageConfig blockstoreType(String blockstoreType) {
    
    this.blockstoreType = blockstoreType;
    return this;
  }

   /**
   * Get blockstoreType
   * @return blockstoreType
  **/
  @javax.annotation.Nonnull
  public String getBlockstoreType() {
    return blockstoreType;
  }


  public void setBlockstoreType(String blockstoreType) {
    this.blockstoreType = blockstoreType;
  }


  public StorageConfig blockstoreNamespaceExample(String blockstoreNamespaceExample) {
    
    this.blockstoreNamespaceExample = blockstoreNamespaceExample;
    return this;
  }

   /**
   * Get blockstoreNamespaceExample
   * @return blockstoreNamespaceExample
  **/
  @javax.annotation.Nonnull
  public String getBlockstoreNamespaceExample() {
    return blockstoreNamespaceExample;
  }


  public void setBlockstoreNamespaceExample(String blockstoreNamespaceExample) {
    this.blockstoreNamespaceExample = blockstoreNamespaceExample;
  }


  public StorageConfig blockstoreNamespaceValidityRegex(String blockstoreNamespaceValidityRegex) {
    
    this.blockstoreNamespaceValidityRegex = blockstoreNamespaceValidityRegex;
    return this;
  }

   /**
   * Get blockstoreNamespaceValidityRegex
   * @return blockstoreNamespaceValidityRegex
  **/
  @javax.annotation.Nonnull
  public String getBlockstoreNamespaceValidityRegex() {
    return blockstoreNamespaceValidityRegex;
  }


  public void setBlockstoreNamespaceValidityRegex(String blockstoreNamespaceValidityRegex) {
    this.blockstoreNamespaceValidityRegex = blockstoreNamespaceValidityRegex;
  }


  public StorageConfig defaultNamespacePrefix(String defaultNamespacePrefix) {
    
    this.defaultNamespacePrefix = defaultNamespacePrefix;
    return this;
  }

   /**
   * Get defaultNamespacePrefix
   * @return defaultNamespacePrefix
  **/
  @javax.annotation.Nullable
  public String getDefaultNamespacePrefix() {
    return defaultNamespacePrefix;
  }


  public void setDefaultNamespacePrefix(String defaultNamespacePrefix) {
    this.defaultNamespacePrefix = defaultNamespacePrefix;
  }


  public StorageConfig preSignSupport(Boolean preSignSupport) {
    
    this.preSignSupport = preSignSupport;
    return this;
  }

   /**
   * Get preSignSupport
   * @return preSignSupport
  **/
  @javax.annotation.Nonnull
  public Boolean getPreSignSupport() {
    return preSignSupport;
  }


  public void setPreSignSupport(Boolean preSignSupport) {
    this.preSignSupport = preSignSupport;
  }


  public StorageConfig preSignSupportUi(Boolean preSignSupportUi) {
    
    this.preSignSupportUi = preSignSupportUi;
    return this;
  }

   /**
   * Get preSignSupportUi
   * @return preSignSupportUi
  **/
  @javax.annotation.Nonnull
  public Boolean getPreSignSupportUi() {
    return preSignSupportUi;
  }


  public void setPreSignSupportUi(Boolean preSignSupportUi) {
    this.preSignSupportUi = preSignSupportUi;
  }


  public StorageConfig importSupport(Boolean importSupport) {
    
    this.importSupport = importSupport;
    return this;
  }

   /**
   * Get importSupport
   * @return importSupport
  **/
  @javax.annotation.Nonnull
  public Boolean getImportSupport() {
    return importSupport;
  }


  public void setImportSupport(Boolean importSupport) {
    this.importSupport = importSupport;
  }


  public StorageConfig importValidityRegex(String importValidityRegex) {
    
    this.importValidityRegex = importValidityRegex;
    return this;
  }

   /**
   * Get importValidityRegex
   * @return importValidityRegex
  **/
  @javax.annotation.Nonnull
  public String getImportValidityRegex() {
    return importValidityRegex;
  }


  public void setImportValidityRegex(String importValidityRegex) {
    this.importValidityRegex = importValidityRegex;
  }


  public StorageConfig preSignMultipartUpload(Boolean preSignMultipartUpload) {
    
    this.preSignMultipartUpload = preSignMultipartUpload;
    return this;
  }

   /**
   * Get preSignMultipartUpload
   * @return preSignMultipartUpload
  **/
  @javax.annotation.Nullable
  public Boolean getPreSignMultipartUpload() {
    return preSignMultipartUpload;
  }


  public void setPreSignMultipartUpload(Boolean preSignMultipartUpload) {
    this.preSignMultipartUpload = preSignMultipartUpload;
  }


  public StorageConfig blockstoreId(String blockstoreId) {
    
    this.blockstoreId = blockstoreId;
    return this;
  }

   /**
   * Get blockstoreId
   * @return blockstoreId
  **/
  @javax.annotation.Nullable
  public String getBlockstoreId() {
    return blockstoreId;
  }


  public void setBlockstoreId(String blockstoreId) {
    this.blockstoreId = blockstoreId;
  }


  public StorageConfig blockstoreDescription(String blockstoreDescription) {
    
    this.blockstoreDescription = blockstoreDescription;
    return this;
  }

   /**
   * Get blockstoreDescription
   * @return blockstoreDescription
  **/
  @javax.annotation.Nullable
  public String getBlockstoreDescription() {
    return blockstoreDescription;
  }


  public void setBlockstoreDescription(String blockstoreDescription) {
    this.blockstoreDescription = blockstoreDescription;
  }

  /**
   * A container for additional, undeclared properties.
   * This is a holder for any undeclared properties as specified with
   * the 'additionalProperties' keyword in the OAS document.
   */
  private Map<String, Object> additionalProperties;

  /**
   * Set the additional (undeclared) property with the specified name and value.
   * If the property does not already exist, create it otherwise replace it.
   *
   * @param key name of the property
   * @param value value of the property
   * @return the StorageConfig instance itself
   */
  public StorageConfig putAdditionalProperty(String key, Object value) {
    if (this.additionalProperties == null) {
        this.additionalProperties = new HashMap<String, Object>();
    }
    this.additionalProperties.put(key, value);
    return this;
  }

  /**
   * Return the additional (undeclared) property.
   *
   * @return a map of objects
   */
  public Map<String, Object> getAdditionalProperties() {
    return additionalProperties;
  }

  /**
   * Return the additional (undeclared) property with the specified name.
   *
   * @param key name of the property
   * @return an object
   */
  public Object getAdditionalProperty(String key) {
    if (this.additionalProperties == null) {
        return null;
    }
    return this.additionalProperties.get(key);
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    StorageConfig storageConfig = (StorageConfig) o;
    return Objects.equals(this.blockstoreType, storageConfig.blockstoreType) &&
        Objects.equals(this.blockstoreNamespaceExample, storageConfig.blockstoreNamespaceExample) &&
        Objects.equals(this.blockstoreNamespaceValidityRegex, storageConfig.blockstoreNamespaceValidityRegex) &&
        Objects.equals(this.defaultNamespacePrefix, storageConfig.defaultNamespacePrefix) &&
        Objects.equals(this.preSignSupport, storageConfig.preSignSupport) &&
        Objects.equals(this.preSignSupportUi, storageConfig.preSignSupportUi) &&
        Objects.equals(this.importSupport, storageConfig.importSupport) &&
        Objects.equals(this.importValidityRegex, storageConfig.importValidityRegex) &&
        Objects.equals(this.preSignMultipartUpload, storageConfig.preSignMultipartUpload) &&
        Objects.equals(this.blockstoreId, storageConfig.blockstoreId) &&
        Objects.equals(this.blockstoreDescription, storageConfig.blockstoreDescription)&&
        Objects.equals(this.additionalProperties, storageConfig.additionalProperties);
  }

  @Override
  public int hashCode() {
    return Objects.hash(blockstoreType, blockstoreNamespaceExample, blockstoreNamespaceValidityRegex, defaultNamespacePrefix, preSignSupport, preSignSupportUi, importSupport, importValidityRegex, preSignMultipartUpload, blockstoreId, blockstoreDescription, additionalProperties);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class StorageConfig {\n");
    sb.append("    blockstoreType: ").append(toIndentedString(blockstoreType)).append("\n");
    sb.append("    blockstoreNamespaceExample: ").append(toIndentedString(blockstoreNamespaceExample)).append("\n");
    sb.append("    blockstoreNamespaceValidityRegex: ").append(toIndentedString(blockstoreNamespaceValidityRegex)).append("\n");
    sb.append("    defaultNamespacePrefix: ").append(toIndentedString(defaultNamespacePrefix)).append("\n");
    sb.append("    preSignSupport: ").append(toIndentedString(preSignSupport)).append("\n");
    sb.append("    preSignSupportUi: ").append(toIndentedString(preSignSupportUi)).append("\n");
    sb.append("    importSupport: ").append(toIndentedString(importSupport)).append("\n");
    sb.append("    importValidityRegex: ").append(toIndentedString(importValidityRegex)).append("\n");
    sb.append("    preSignMultipartUpload: ").append(toIndentedString(preSignMultipartUpload)).append("\n");
    sb.append("    blockstoreId: ").append(toIndentedString(blockstoreId)).append("\n");
    sb.append("    blockstoreDescription: ").append(toIndentedString(blockstoreDescription)).append("\n");
    sb.append("    additionalProperties: ").append(toIndentedString(additionalProperties)).append("\n");
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


  public static HashSet<String> openapiFields;
  public static HashSet<String> openapiRequiredFields;

  static {
    // a set of all properties/fields (JSON key names)
    openapiFields = new HashSet<String>();
    openapiFields.add("blockstore_type");
    openapiFields.add("blockstore_namespace_example");
    openapiFields.add("blockstore_namespace_ValidityRegex");
    openapiFields.add("default_namespace_prefix");
    openapiFields.add("pre_sign_support");
    openapiFields.add("pre_sign_support_ui");
    openapiFields.add("import_support");
    openapiFields.add("import_validity_regex");
    openapiFields.add("pre_sign_multipart_upload");
    openapiFields.add("blockstore_id");
    openapiFields.add("blockstore_description");

    // a set of required properties/fields (JSON key names)
    openapiRequiredFields = new HashSet<String>();
    openapiRequiredFields.add("blockstore_type");
    openapiRequiredFields.add("blockstore_namespace_example");
    openapiRequiredFields.add("blockstore_namespace_ValidityRegex");
    openapiRequiredFields.add("pre_sign_support");
    openapiRequiredFields.add("pre_sign_support_ui");
    openapiRequiredFields.add("import_support");
    openapiRequiredFields.add("import_validity_regex");
  }

 /**
  * Validates the JSON Element and throws an exception if issues found
  *
  * @param jsonElement JSON Element
  * @throws IOException if the JSON Element is invalid with respect to StorageConfig
  */
  public static void validateJsonElement(JsonElement jsonElement) throws IOException {
      if (jsonElement == null) {
        if (!StorageConfig.openapiRequiredFields.isEmpty()) { // has required fields but JSON element is null
          throw new IllegalArgumentException(String.format("The required field(s) %s in StorageConfig is not found in the empty JSON string", StorageConfig.openapiRequiredFields.toString()));
        }
      }

      // check to make sure all required properties/fields are present in the JSON string
      for (String requiredField : StorageConfig.openapiRequiredFields) {
        if (jsonElement.getAsJsonObject().get(requiredField) == null) {
          throw new IllegalArgumentException(String.format("The required field `%s` is not found in the JSON string: %s", requiredField, jsonElement.toString()));
        }
      }
        JsonObject jsonObj = jsonElement.getAsJsonObject();
      if (!jsonObj.get("blockstore_type").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `blockstore_type` to be a primitive type in the JSON string but got `%s`", jsonObj.get("blockstore_type").toString()));
      }
      if (!jsonObj.get("blockstore_namespace_example").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `blockstore_namespace_example` to be a primitive type in the JSON string but got `%s`", jsonObj.get("blockstore_namespace_example").toString()));
      }
      if (!jsonObj.get("blockstore_namespace_ValidityRegex").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `blockstore_namespace_ValidityRegex` to be a primitive type in the JSON string but got `%s`", jsonObj.get("blockstore_namespace_ValidityRegex").toString()));
      }
      if ((jsonObj.get("default_namespace_prefix") != null && !jsonObj.get("default_namespace_prefix").isJsonNull()) && !jsonObj.get("default_namespace_prefix").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `default_namespace_prefix` to be a primitive type in the JSON string but got `%s`", jsonObj.get("default_namespace_prefix").toString()));
      }
      if (!jsonObj.get("import_validity_regex").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `import_validity_regex` to be a primitive type in the JSON string but got `%s`", jsonObj.get("import_validity_regex").toString()));
      }
      if ((jsonObj.get("blockstore_id") != null && !jsonObj.get("blockstore_id").isJsonNull()) && !jsonObj.get("blockstore_id").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `blockstore_id` to be a primitive type in the JSON string but got `%s`", jsonObj.get("blockstore_id").toString()));
      }
      if ((jsonObj.get("blockstore_description") != null && !jsonObj.get("blockstore_description").isJsonNull()) && !jsonObj.get("blockstore_description").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `blockstore_description` to be a primitive type in the JSON string but got `%s`", jsonObj.get("blockstore_description").toString()));
      }
  }

  public static class CustomTypeAdapterFactory implements TypeAdapterFactory {
    @SuppressWarnings("unchecked")
    @Override
    public <T> TypeAdapter<T> create(Gson gson, TypeToken<T> type) {
       if (!StorageConfig.class.isAssignableFrom(type.getRawType())) {
         return null; // this class only serializes 'StorageConfig' and its subtypes
       }
       final TypeAdapter<JsonElement> elementAdapter = gson.getAdapter(JsonElement.class);
       final TypeAdapter<StorageConfig> thisAdapter
                        = gson.getDelegateAdapter(this, TypeToken.get(StorageConfig.class));

       return (TypeAdapter<T>) new TypeAdapter<StorageConfig>() {
           @Override
           public void write(JsonWriter out, StorageConfig value) throws IOException {
             JsonObject obj = thisAdapter.toJsonTree(value).getAsJsonObject();
             obj.remove("additionalProperties");
             // serialize additional properties
             if (value.getAdditionalProperties() != null) {
               for (Map.Entry<String, Object> entry : value.getAdditionalProperties().entrySet()) {
                 if (entry.getValue() instanceof String)
                   obj.addProperty(entry.getKey(), (String) entry.getValue());
                 else if (entry.getValue() instanceof Number)
                   obj.addProperty(entry.getKey(), (Number) entry.getValue());
                 else if (entry.getValue() instanceof Boolean)
                   obj.addProperty(entry.getKey(), (Boolean) entry.getValue());
                 else if (entry.getValue() instanceof Character)
                   obj.addProperty(entry.getKey(), (Character) entry.getValue());
                 else {
                   obj.add(entry.getKey(), gson.toJsonTree(entry.getValue()).getAsJsonObject());
                 }
               }
             }
             elementAdapter.write(out, obj);
           }

           @Override
           public StorageConfig read(JsonReader in) throws IOException {
             JsonElement jsonElement = elementAdapter.read(in);
             validateJsonElement(jsonElement);
             JsonObject jsonObj = jsonElement.getAsJsonObject();
             // store additional fields in the deserialized instance
             StorageConfig instance = thisAdapter.fromJsonTree(jsonObj);
             for (Map.Entry<String, JsonElement> entry : jsonObj.entrySet()) {
               if (!openapiFields.contains(entry.getKey())) {
                 if (entry.getValue().isJsonPrimitive()) { // primitive type
                   if (entry.getValue().getAsJsonPrimitive().isString())
                     instance.putAdditionalProperty(entry.getKey(), entry.getValue().getAsString());
                   else if (entry.getValue().getAsJsonPrimitive().isNumber())
                     instance.putAdditionalProperty(entry.getKey(), entry.getValue().getAsNumber());
                   else if (entry.getValue().getAsJsonPrimitive().isBoolean())
                     instance.putAdditionalProperty(entry.getKey(), entry.getValue().getAsBoolean());
                   else
                     throw new IllegalArgumentException(String.format("The field `%s` has unknown primitive type. Value: %s", entry.getKey(), entry.getValue().toString()));
                 } else if (entry.getValue().isJsonArray()) {
                     instance.putAdditionalProperty(entry.getKey(), gson.fromJson(entry.getValue(), List.class));
                 } else { // JSON object
                     instance.putAdditionalProperty(entry.getKey(), gson.fromJson(entry.getValue(), HashMap.class));
                 }
               }
             }
             return instance;
           }

       }.nullSafe();
    }
  }

 /**
  * Create an instance of StorageConfig given an JSON string
  *
  * @param jsonString JSON string
  * @return An instance of StorageConfig
  * @throws IOException if the JSON string is invalid with respect to StorageConfig
  */
  public static StorageConfig fromJson(String jsonString) throws IOException {
    return JSON.getGson().fromJson(jsonString, StorageConfig.class);
  }

 /**
  * Convert an instance of StorageConfig to an JSON string
  *
  * @return JSON string
  */
  public String toJson() {
    return JSON.getGson().toJson(this);
  }
}

