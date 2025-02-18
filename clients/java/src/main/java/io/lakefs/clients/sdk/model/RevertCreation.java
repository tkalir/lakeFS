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
import io.lakefs.clients.sdk.model.CommitOverrides;
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
 * RevertCreation
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class RevertCreation {
  public static final String SERIALIZED_NAME_REF = "ref";
  @SerializedName(SERIALIZED_NAME_REF)
  private String ref;

  public static final String SERIALIZED_NAME_COMMIT_OVERRIDES = "commit_overrides";
  @SerializedName(SERIALIZED_NAME_COMMIT_OVERRIDES)
  private CommitOverrides commitOverrides;

  public static final String SERIALIZED_NAME_PARENT_NUMBER = "parent_number";
  @SerializedName(SERIALIZED_NAME_PARENT_NUMBER)
  private Integer parentNumber;

  public static final String SERIALIZED_NAME_FORCE = "force";
  @SerializedName(SERIALIZED_NAME_FORCE)
  private Boolean force = false;

  public static final String SERIALIZED_NAME_ALLOW_EMPTY = "allow_empty";
  @SerializedName(SERIALIZED_NAME_ALLOW_EMPTY)
  private Boolean allowEmpty = false;

  public RevertCreation() {
  }

  public RevertCreation ref(String ref) {
    
    this.ref = ref;
    return this;
  }

   /**
   * the commit to revert, given by a ref
   * @return ref
  **/
  @javax.annotation.Nonnull
  public String getRef() {
    return ref;
  }


  public void setRef(String ref) {
    this.ref = ref;
  }


  public RevertCreation commitOverrides(CommitOverrides commitOverrides) {
    
    this.commitOverrides = commitOverrides;
    return this;
  }

   /**
   * Get commitOverrides
   * @return commitOverrides
  **/
  @javax.annotation.Nullable
  public CommitOverrides getCommitOverrides() {
    return commitOverrides;
  }


  public void setCommitOverrides(CommitOverrides commitOverrides) {
    this.commitOverrides = commitOverrides;
  }


  public RevertCreation parentNumber(Integer parentNumber) {
    
    this.parentNumber = parentNumber;
    return this;
  }

   /**
   * when reverting a merge commit, the parent number (starting from 1) relative to which to perform the revert.
   * @return parentNumber
  **/
  @javax.annotation.Nonnull
  public Integer getParentNumber() {
    return parentNumber;
  }


  public void setParentNumber(Integer parentNumber) {
    this.parentNumber = parentNumber;
  }


  public RevertCreation force(Boolean force) {
    
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


  public RevertCreation allowEmpty(Boolean allowEmpty) {
    
    this.allowEmpty = allowEmpty;
    return this;
  }

   /**
   * allow empty commit (revert without changes)
   * @return allowEmpty
  **/
  @javax.annotation.Nullable
  public Boolean getAllowEmpty() {
    return allowEmpty;
  }


  public void setAllowEmpty(Boolean allowEmpty) {
    this.allowEmpty = allowEmpty;
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
   * @return the RevertCreation instance itself
   */
  public RevertCreation putAdditionalProperty(String key, Object value) {
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
    RevertCreation revertCreation = (RevertCreation) o;
    return Objects.equals(this.ref, revertCreation.ref) &&
        Objects.equals(this.commitOverrides, revertCreation.commitOverrides) &&
        Objects.equals(this.parentNumber, revertCreation.parentNumber) &&
        Objects.equals(this.force, revertCreation.force) &&
        Objects.equals(this.allowEmpty, revertCreation.allowEmpty)&&
        Objects.equals(this.additionalProperties, revertCreation.additionalProperties);
  }

  @Override
  public int hashCode() {
    return Objects.hash(ref, commitOverrides, parentNumber, force, allowEmpty, additionalProperties);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class RevertCreation {\n");
    sb.append("    ref: ").append(toIndentedString(ref)).append("\n");
    sb.append("    commitOverrides: ").append(toIndentedString(commitOverrides)).append("\n");
    sb.append("    parentNumber: ").append(toIndentedString(parentNumber)).append("\n");
    sb.append("    force: ").append(toIndentedString(force)).append("\n");
    sb.append("    allowEmpty: ").append(toIndentedString(allowEmpty)).append("\n");
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
    openapiFields.add("ref");
    openapiFields.add("commit_overrides");
    openapiFields.add("parent_number");
    openapiFields.add("force");
    openapiFields.add("allow_empty");

    // a set of required properties/fields (JSON key names)
    openapiRequiredFields = new HashSet<String>();
    openapiRequiredFields.add("ref");
    openapiRequiredFields.add("parent_number");
  }

 /**
  * Validates the JSON Element and throws an exception if issues found
  *
  * @param jsonElement JSON Element
  * @throws IOException if the JSON Element is invalid with respect to RevertCreation
  */
  public static void validateJsonElement(JsonElement jsonElement) throws IOException {
      if (jsonElement == null) {
        if (!RevertCreation.openapiRequiredFields.isEmpty()) { // has required fields but JSON element is null
          throw new IllegalArgumentException(String.format("The required field(s) %s in RevertCreation is not found in the empty JSON string", RevertCreation.openapiRequiredFields.toString()));
        }
      }

      // check to make sure all required properties/fields are present in the JSON string
      for (String requiredField : RevertCreation.openapiRequiredFields) {
        if (jsonElement.getAsJsonObject().get(requiredField) == null) {
          throw new IllegalArgumentException(String.format("The required field `%s` is not found in the JSON string: %s", requiredField, jsonElement.toString()));
        }
      }
        JsonObject jsonObj = jsonElement.getAsJsonObject();
      if (!jsonObj.get("ref").isJsonPrimitive()) {
        throw new IllegalArgumentException(String.format("Expected the field `ref` to be a primitive type in the JSON string but got `%s`", jsonObj.get("ref").toString()));
      }
      // validate the optional field `commit_overrides`
      if (jsonObj.get("commit_overrides") != null && !jsonObj.get("commit_overrides").isJsonNull()) {
        CommitOverrides.validateJsonElement(jsonObj.get("commit_overrides"));
      }
  }

  public static class CustomTypeAdapterFactory implements TypeAdapterFactory {
    @SuppressWarnings("unchecked")
    @Override
    public <T> TypeAdapter<T> create(Gson gson, TypeToken<T> type) {
       if (!RevertCreation.class.isAssignableFrom(type.getRawType())) {
         return null; // this class only serializes 'RevertCreation' and its subtypes
       }
       final TypeAdapter<JsonElement> elementAdapter = gson.getAdapter(JsonElement.class);
       final TypeAdapter<RevertCreation> thisAdapter
                        = gson.getDelegateAdapter(this, TypeToken.get(RevertCreation.class));

       return (TypeAdapter<T>) new TypeAdapter<RevertCreation>() {
           @Override
           public void write(JsonWriter out, RevertCreation value) throws IOException {
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
           public RevertCreation read(JsonReader in) throws IOException {
             JsonElement jsonElement = elementAdapter.read(in);
             validateJsonElement(jsonElement);
             JsonObject jsonObj = jsonElement.getAsJsonObject();
             // store additional fields in the deserialized instance
             RevertCreation instance = thisAdapter.fromJsonTree(jsonObj);
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
  * Create an instance of RevertCreation given an JSON string
  *
  * @param jsonString JSON string
  * @return An instance of RevertCreation
  * @throws IOException if the JSON string is invalid with respect to RevertCreation
  */
  public static RevertCreation fromJson(String jsonString) throws IOException {
    return JSON.getGson().fromJson(jsonString, RevertCreation.class);
  }

 /**
  * Convert an instance of RevertCreation to an JSON string
  *
  * @return JSON string
  */
  public String toJson() {
    return JSON.getGson().toJson(this);
  }
}

