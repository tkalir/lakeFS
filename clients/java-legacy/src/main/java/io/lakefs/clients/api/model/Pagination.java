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
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * Pagination
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class Pagination {
  public static final String SERIALIZED_NAME_HAS_MORE = "has_more";
  @SerializedName(SERIALIZED_NAME_HAS_MORE)
  private Boolean hasMore;

  public static final String SERIALIZED_NAME_NEXT_OFFSET = "next_offset";
  @SerializedName(SERIALIZED_NAME_NEXT_OFFSET)
  private String nextOffset;

  public static final String SERIALIZED_NAME_RESULTS = "results";
  @SerializedName(SERIALIZED_NAME_RESULTS)
  private Integer results;

  public static final String SERIALIZED_NAME_MAX_PER_PAGE = "max_per_page";
  @SerializedName(SERIALIZED_NAME_MAX_PER_PAGE)
  private Integer maxPerPage;


  public Pagination hasMore(Boolean hasMore) {
    
    this.hasMore = hasMore;
    return this;
  }

   /**
   * Next page is available
   * @return hasMore
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "Next page is available")

  public Boolean getHasMore() {
    return hasMore;
  }


  public void setHasMore(Boolean hasMore) {
    this.hasMore = hasMore;
  }


  public Pagination nextOffset(String nextOffset) {
    
    this.nextOffset = nextOffset;
    return this;
  }

   /**
   * Token used to retrieve the next page
   * @return nextOffset
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "Token used to retrieve the next page")

  public String getNextOffset() {
    return nextOffset;
  }


  public void setNextOffset(String nextOffset) {
    this.nextOffset = nextOffset;
  }


  public Pagination results(Integer results) {
    
    this.results = results;
    return this;
  }

   /**
   * Number of values found in the results
   * minimum: 0
   * @return results
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "Number of values found in the results")

  public Integer getResults() {
    return results;
  }


  public void setResults(Integer results) {
    this.results = results;
  }


  public Pagination maxPerPage(Integer maxPerPage) {
    
    this.maxPerPage = maxPerPage;
    return this;
  }

   /**
   * Maximal number of entries per page
   * minimum: 0
   * @return maxPerPage
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "Maximal number of entries per page")

  public Integer getMaxPerPage() {
    return maxPerPage;
  }


  public void setMaxPerPage(Integer maxPerPage) {
    this.maxPerPage = maxPerPage;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    Pagination pagination = (Pagination) o;
    return Objects.equals(this.hasMore, pagination.hasMore) &&
        Objects.equals(this.nextOffset, pagination.nextOffset) &&
        Objects.equals(this.results, pagination.results) &&
        Objects.equals(this.maxPerPage, pagination.maxPerPage);
  }

  @Override
  public int hashCode() {
    return Objects.hash(hasMore, nextOffset, results, maxPerPage);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class Pagination {\n");
    sb.append("    hasMore: ").append(toIndentedString(hasMore)).append("\n");
    sb.append("    nextOffset: ").append(toIndentedString(nextOffset)).append("\n");
    sb.append("    results: ").append(toIndentedString(results)).append("\n");
    sb.append("    maxPerPage: ").append(toIndentedString(maxPerPage)).append("\n");
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

