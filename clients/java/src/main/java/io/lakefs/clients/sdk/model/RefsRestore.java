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
 * RefsRestore
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class RefsRestore {
  public static final String SERIALIZED_NAME_COMMITS_META_RANGE_ID = "commits_meta_range_id";
  @SerializedName(SERIALIZED_NAME_COMMITS_META_RANGE_ID)
  private String commitsMetaRangeId;

  public static final String SERIALIZED_NAME_TAGS_META_RANGE_ID = "tags_meta_range_id";
  @SerializedName(SERIALIZED_NAME_TAGS_META_RANGE_ID)
  private String tagsMetaRangeId;

  public static final String SERIALIZED_NAME_BRANCHES_META_RANGE_ID = "branches_meta_range_id";
  @SerializedName(SERIALIZED_NAME_BRANCHES_META_RANGE_ID)
  private String branchesMetaRangeId;

  public static final String SERIALIZED_NAME_FORCE = "force";
  @SerializedName(SERIALIZED_NAME_FORCE)
  private Boolean force = false;

  public RefsRestore() {
  }

  public RefsRestore commitsMetaRangeId(String commitsMetaRangeId) {
    
    this.commitsMetaRangeId = commitsMetaRangeId;
    return this;
  }

   /**
   * Get commitsMetaRangeId
   * @return commitsMetaRangeId
  **/
  @javax.annotation.Nonnull
  public String getCommitsMetaRangeId() {
    return commitsMetaRangeId;
  }


  public void setCommitsMetaRangeId(String commitsMetaRangeId) {
    this.commitsMetaRangeId = commitsMetaRangeId;
  }


  public RefsRestore tagsMetaRangeId(String tagsMetaRangeId) {
    
    this.tagsMetaRangeId = tagsMetaRangeId;
    return this;
  }

   /**
   * Get tagsMetaRangeId
   * @return tagsMetaRangeId
  **/
  @javax.annotation.Nonnull
  public String getTagsMetaRangeId() {
    return tagsMetaRangeId;
  }


  public void setTagsMetaRangeId(String tagsMetaRangeId) {
    this.tagsMetaRangeId = tagsMetaRangeId;
  }


  public RefsRestore branchesMetaRangeId(String branchesMetaRangeId) {
    
    this.branchesMetaRangeId = branchesMetaRangeId;
    return this;
  }

   /**
   * Get branchesMetaRangeId
   * @return branchesMetaRangeId
  **/
  @javax.annotation.Nonnull
  public String getBranchesMetaRangeId() {
    return branchesMetaRangeId;
  }


  public void setBranchesMetaRangeId(String branchesMetaRangeId) {
    this.branchesMetaRangeId = branchesMetaRangeId;
  }


  public RefsRestore force(Boolean force) {
    
    this.force = force;
    return this;
  }

   /**
   * Get force
   * @return force
  **/
  @javax.annotation.Nullable
  public Boolean getForce() {
    return force;
  }


  public void setForce(Boolean force) {
    this.force = force;
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
   * @return the RefsRestore instance itself
   */
  public RefsRestore putAdditionalProperty(String key, Object value) {
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
    RefsRestore refsRestore = (RefsRestore) o;
    return Objects.equals(this.commitsMetaRangeId, refsRestore.commitsMetaRangeId) &&
        Objects.equals(this.tagsMetaRangeId, refsRestore.tagsMetaRangeId) &&
        Objects.equals(this.branchesMetaRangeId, refsRestore.branchesMetaRangeId) &&
        Objects.equals(this.force, refsRestore.force)&&
        Objects.equals(this.additionalProperties, refsRestore.additionalProperties);
  }

  @Override
  public int hashCode() {
    return Objects.hash(commitsMetaRangeId, tagsMetaRangeId, branchesMetaRangeId, force, additionalProperties);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class RefsRestore {\n");
    sb.append("    commitsMetaRangeId: ").append(toIndentedString(commitsMetaRangeId)).append("\n");
    sb.append("    tagsMetaRangeId: ").append(toIndentedString(tagsMetaRangeId)).append("\n");
    sb.append("    branchesMetaRangeId: ").append(toIndentedString(branchesMetaRangeId)).append("\n");
    sb.append("    force: ").append(toIndentedString(force)).append("\n");
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
    openapiFields.add("commits_meta_range_id");
    openapiFields.add("tags_meta_range_id");
    openapiFields.add("branches_meta_range_id");
    openapiFields.add("force");

    // a set of required properties/fields (JSON key names)
    openapiRequiredFields = new HashSet<String>();
    openapiRequiredFields.add("commits_meta_range_id");
    openapiRequiredFields.add("tags_meta_range_id");
    openapiRequiredFields.add("branches_meta_range_id");
  }

 /**
  * Validates the JSON Element and throws an exception if issues found
  *
  * @param jsonElement JSON Element
  * @throws IOException if the JSON Element is invalid with respect to RefsRestore
  */
  public static void validateJsonElement(JsonElement jsonElement) throws IOException {
      if (jsonElement == null) {
        if (!RefsRestore.openapiRequiredFields.isEmpty()) { // has required fields but JSON element is null
          throw new IllegalArgumentException(String.format("The required field(s) %s in RefsRestore is not found in the empty JSON string", RefsRestore.openapiRequiredFields.toString()));
        }
      }

      // check to make sure all required properties/fields are present in the JSON string
      for (String requiredField : RefsRestore.openapiRequiredFields) {
        if (jsonElement.getAsJsonObject().get(requiredField) == null) {
          throw new IllegalArgumentException(String.format("The required field `%s` is not found in the JSON string: %s", requiredField, jsonElement.toString()));
        }
      }
        JsonObject jsonObj = jsonElement.getAsJsonObject();
      if (!jsonObj.get("commits_meta_range_id").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `commits_meta_range_id` to be a primitive type in the JSON string but got `%s`", jsonObj.get("commits_meta_range_id").toString()));
      }
      if (!jsonObj.get("tags_meta_range_id").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `tags_meta_range_id` to be a primitive type in the JSON string but got `%s`", jsonObj.get("tags_meta_range_id").toString()));
      }
      if (!jsonObj.get("branches_meta_range_id").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `branches_meta_range_id` to be a primitive type in the JSON string but got `%s`", jsonObj.get("branches_meta_range_id").toString()));
      }
  }

  public static class CustomTypeAdapterFactory implements TypeAdapterFactory {
    @SuppressWarnings("unchecked")
    @Override
    public <T> TypeAdapter<T> create(Gson gson, TypeToken<T> type) {
       if (!RefsRestore.class.isAssignableFrom(type.getRawType())) {
         return null; // this class only serializes 'RefsRestore' and its subtypes
       }
       final TypeAdapter<JsonElement> elementAdapter = gson.getAdapter(JsonElement.class);
       final TypeAdapter<RefsRestore> thisAdapter
                        = gson.getDelegateAdapter(this, TypeToken.get(RefsRestore.class));

       return (TypeAdapter<T>) new TypeAdapter<RefsRestore>() {
           @Override
           public void write(JsonWriter out, RefsRestore value) throws IOException {
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
           public RefsRestore read(JsonReader in) throws IOException {
             JsonElement jsonElement = elementAdapter.read(in);
             validateJsonElement(jsonElement);
             JsonObject jsonObj = jsonElement.getAsJsonObject();
             // store additional fields in the deserialized instance
             RefsRestore instance = thisAdapter.fromJsonTree(jsonObj);
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
  * Create an instance of RefsRestore given an JSON string
  *
  * @param jsonString JSON string
  * @return An instance of RefsRestore
  * @throws IOException if the JSON string is invalid with respect to RefsRestore
  */
  public static RefsRestore fromJson(String jsonString) throws IOException {
    return JSON.getGson().fromJson(jsonString, RefsRestore.class);
  }

 /**
  * Convert an instance of RefsRestore to an JSON string
  *
  * @return JSON string
  */
  public String toJson() {
    return JSON.getGson().toJson(this);
  }
}

