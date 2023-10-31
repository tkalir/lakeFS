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


package io.lakefs.clients.api.model;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.lakefs.clients.api.model.OtfDiffEntry;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * OtfDiffList
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class OtfDiffList {
  /**
   * Gets or Sets diffType
   */
  @JsonAdapter(DiffTypeEnum.Adapter.class)
  public enum DiffTypeEnum {
    CREATED("created"),
    
    DROPPED("dropped"),
    
    CHANGED("changed");

    private String value;

    DiffTypeEnum(String value) {
      this.value = value;
    }

    public String getValue() {
      return value;
    }

    @Override
    public String toString() {
      return String.valueOf(value);
    }

    public static DiffTypeEnum fromValue(String value) {
      for (DiffTypeEnum b : DiffTypeEnum.values()) {
        if (b.value.equals(value)) {
          return b;
        }
      }
      throw new IllegalArgumentException("Unexpected value '" + value + "'");
    }

    public static class Adapter extends TypeAdapter<DiffTypeEnum> {
      @Override
      public void write(final JsonWriter jsonWriter, final DiffTypeEnum enumeration) throws IOException {
        jsonWriter.value(enumeration.getValue());
      }

      @Override
      public DiffTypeEnum read(final JsonReader jsonReader) throws IOException {
        String value =  jsonReader.nextString();
        return DiffTypeEnum.fromValue(value);
      }
    }
  }

  public static final String SERIALIZED_NAME_DIFF_TYPE = "diff_type";
  @SerializedName(SERIALIZED_NAME_DIFF_TYPE)
  private DiffTypeEnum diffType;

  public static final String SERIALIZED_NAME_RESULTS = "results";
  @SerializedName(SERIALIZED_NAME_RESULTS)
  private List<OtfDiffEntry> results = new ArrayList<OtfDiffEntry>();


  public OtfDiffList diffType(DiffTypeEnum diffType) {
    
    this.diffType = diffType;
    return this;
  }

   /**
   * Get diffType
   * @return diffType
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public DiffTypeEnum getDiffType() {
    return diffType;
  }


  public void setDiffType(DiffTypeEnum diffType) {
    this.diffType = diffType;
  }


  public OtfDiffList results(List<OtfDiffEntry> results) {
    
    this.results = results;
    return this;
  }

  public OtfDiffList addResultsItem(OtfDiffEntry resultsItem) {
    this.results.add(resultsItem);
    return this;
  }

   /**
   * Get results
   * @return results
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public List<OtfDiffEntry> getResults() {
    return results;
  }


  public void setResults(List<OtfDiffEntry> results) {
    this.results = results;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    OtfDiffList otfDiffList = (OtfDiffList) o;
    return Objects.equals(this.diffType, otfDiffList.diffType) &&
        Objects.equals(this.results, otfDiffList.results);
  }

  @Override
  public int hashCode() {
    return Objects.hash(diffType, results);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class OtfDiffList {\n");
    sb.append("    diffType: ").append(toIndentedString(diffType)).append("\n");
    sb.append("    results: ").append(toIndentedString(results)).append("\n");
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

