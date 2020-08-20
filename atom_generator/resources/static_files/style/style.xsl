<?xml version="1.0" encoding="utf-8" ?>
<xsl:stylesheet version="2.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns:atom="http://www.w3.org/2005/Atom"
                xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xs="http://www.w3.org/2001/XMLSchema"
                xmlns:georss="http://www.georss.org/georss">
  <xsl:output method="html" indent="yes" omit-xml-declaration="no" encoding="UTF-8"/>
  <xsl:template match="/">
    <xsl:variable name="isParent">
      <xsl:choose>
        <xsl:when test="atom:feed/atom:entry[1]/atom:link[@type='application/atom+xml']/@href">
          <xsl:value-of select="'true'"></xsl:value-of>
        </xsl:when>
        <xsl:otherwise>
          <xsl:value-of select="'false'"></xsl:value-of>
        </xsl:otherwise>
      </xsl:choose>

    </xsl:variable>
    <html>
      <head>
        <title>
          <xsl:value-of select="atom:feed/atom:title"/>
        </title>
        <link rel="stylesheet" type="text/css" href="./assets/deps.css"/>
        <link rel="stylesheet" type="text/css" href="./style/style.css"/>
        <script src="./assets/deps.js"/>
      </head>

      <body>
        <header id="banner" role="banner">
          <div id="heading">
            <nav role="navigation" class="full-width navbar navbar-default">
              <div class="content">
                <a class="navbar-brand" href="https://www.pdok.nl">
                  <img src="./assets/images/pdok-logo.png" alt="Logo PDOK: Ga naar de homepage"/>
                </a>
              </div>
            </nav>
          </div>
        </header>

        <div class="content">
          <div id="feedTitleContainer" class="entry">
            <h1 id="feedTitleText">
              <xsl:choose>
                <xsl:when test="$isParent='true'">
                  <xsl:text>Service Feed -</xsl:text>
                </xsl:when>
                <xsl:otherwise>
                  <xsl:text>Data Feed -</xsl:text>
                </xsl:otherwise>
              </xsl:choose>
              <xsl:value-of select="atom:feed/atom:title"/>
            </h1>
            <h2 id="feedSubtitleText">
              <xsl:value-of select="atom:feed/atom:subtitle"/>
            </h2>
            <table>
              <xsl:if test="$isParent = 'false'">
                <tr>
                  <td>Main page</td>
                  <td>
                    <a href="{atom:feed/atom:link[@rel='up']/@href}">
                      Parent
                    </a>
                  </td>
                </tr>
              </xsl:if>
              <xsl:if test="$isParent = 'true'">
                <xsl:if test="atom:feed/atom:author">
                  <tr>
                    <td>Service Provider</td>
                    <td>
                      <a href="mailto:{atom:feed/atom:author/atom:email}">
                        <xsl:value-of select="atom:feed/atom:author/atom:name" disable-output-escaping="yes"/>
                      </a>
                    </td>
                  </tr>
                </xsl:if>
              </xsl:if>
              <xsl:if test="$isParent = 'false'">
                <xsl:if test="atom:feed/atom:rights">
                  <tr>
                    <td>Rights</td>
                    <td>
                      <xsl:value-of select="atom:feed/atom:rights" disable-output-escaping="yes"/>
                    </td>
                  </tr>
                </xsl:if>
              </xsl:if>
              <xsl:if test="atom:feed/atom:updated">
                <tr>
                  <td>Updated</td>
                  <td>
                    <xsl:apply-templates select="atom:feed/atom:updated" disable-output-escaping="yes"/>
                  </td>
                </tr>
              </xsl:if>
              <xsl:if test="atom:feed/atom:link[@rel='describedby']/@href">
                <tr>
                  <td>
                    <xsl:choose>
                      <xsl:when test="$isParent='true'">
                        <xsl:text>Service Metadata</xsl:text>
                      </xsl:when>
                      <xsl:otherwise>
                        <xsl:text>Dataset Metadata</xsl:text>
                      </xsl:otherwise>
                    </xsl:choose>
                  </td>
                  <td>
                    <a href="{atom:feed/atom:link[@rel='describedby']/@href}">XML</a>
                    /
                    <a href="{atom:feed/atom:link[@rel='related']/@href}">NGR
                    </a>
                  </td>
                </tr>
              </xsl:if>
              <tr>
                <td>ATOM Feed XML</td>
                <td>
                  <a id="show" title="Show ATOM XML">Show</a>
                </td>
              </tr>
            </table>
            <div id="xml-wrapper" style="display:none;">
              <pre>
                <code class="language-markup" id="atom-xml"></code>
              </pre>
            </div>
          </div>
          <xsl:apply-templates select="//atom:entry">
            <xsl:with-param name="isParent" select="$isParent"></xsl:with-param>
          </xsl:apply-templates>
        </div>
        <script src="./style/script.js"></script>
      </body>

    </html>
  </xsl:template>

  <xsl:template match="atom:updated">
    <xsl:value-of
        select="concat(substring(current(),9,2), '-', substring(current(),6,2), '-', substring(current(),1,4))"/>
  </xsl:template>


  <xsl:template match="atom:entry">
    <xsl:param name="isParent"/>
    <div id="feedContent">
      <div class="entry">
        <h3 id="feedEntryTitle">
          <xsl:choose>
            <xsl:when test="$isParent='true'">
              <a href="{atom:link[@type='application/atom+xml']/@href}">
                <xsl:value-of select="atom:title" disable-output-escaping="yes"/>
              </a>
            </xsl:when>
            <xsl:otherwise>
              <xsl:value-of select="atom:title" disable-output-escaping="yes"/>
            </xsl:otherwise>
          </xsl:choose>
        </h3>

        <xsl:if test="$isParent = 'false'">
          <p>
            <xsl:value-of select="atom:content" disable-output-escaping="yes"/>
          </p>
        </xsl:if>

        <table>
          <xsl:if test="atom:updated">
            <tr>
              <td>Updated</td>
              <td>
                <xsl:apply-templates select="atom:updated" disable-output-escaping="yes"/>
              </td>
            </tr>
          </xsl:if>
          <xsl:if test="georss:polygon">
            <tr>
              <td>Map area</td>

              <td>
                <xsl:variable name="extent">
                  <xsl:apply-templates select="georss:polygon" disable-output-escaping="yes"/>
                </xsl:variable>
                <a href="http://bboxfinder.com/#{$extent}">
                  <xsl:value-of select="$extent" disable-output-escaping="yes"/>
                </a>
              </td>
            </tr>
          </xsl:if>
          <xsl:if test="atom:category/@term">
            <tr>
              <td>Projection</td>
              <td>
                <a href="{atom:category/@term}">
                  <xsl:value-of select="atom:category/@label" disable-output-escaping="yes"/>
                </a>
              </td>
            </tr>
          </xsl:if>
        </table>
        <xsl:choose>
          <xsl:when test="$isParent='false'">
            <xsl:if test="atom:link">
              <a href="{atom:link/@href}" class="download btn btn-default" title="{atom:link/@title}">DOWNLOAD</a>
            </xsl:if>
          </xsl:when>
        </xsl:choose>
      </div>
    </div>
  </xsl:template>

  <xsl:template match="georss:polygon">
    <xsl:variable name="item1" select="substring-before(current(),' ')"/>
    <xsl:variable name="input2" select="substring-after(current(),' ')"/>
    <xsl:variable name="item2" select="substring-before($input2,' ')"/>
    <xsl:variable name="input3" select="substring-after($input2,' ')"/>
    <xsl:variable name="item3" select="substring-before($input3,' ')"/>
    <xsl:variable name="input4" select="substring-after($input3,' ')"/>
    <xsl:variable name="item4" select="substring-before($input4,' ')"/>
    <xsl:variable name="input5" select="substring-after($input4,' ')"/>
    <xsl:variable name="item5" select="substring-before($input5,' ')"/>
    <xsl:variable name="input6" select="substring-after($input5,' ')"/>
    <xsl:variable name="item6" select="substring-before($input6,' ')"/>
    <xsl:value-of select="concat($item1,',',$item2,',',$item5,',',$item6)"/>
  </xsl:template>

  <xsl:template name="string-replace-all">
    <xsl:param name="text"/>
    <xsl:param name="replace"/>
    <xsl:param name="by"/>
    <xsl:choose>
      <xsl:when test="$text = '' or $replace = ''or not($replace)">
        <!-- Prevent this routine from hanging -->
        <xsl:value-of select="$text"/>
      </xsl:when>
      <xsl:when test="contains($text, $replace)">
        <xsl:value-of select="substring-before($text,$replace)"/>
        <xsl:value-of select="$by"/>
        <xsl:call-template name="string-replace-all">
          <xsl:with-param name="text" select="substring-after($text,$replace)"/>
          <xsl:with-param name="replace" select="$replace"/>
          <xsl:with-param name="by" select="$by"/>
        </xsl:call-template>
      </xsl:when>
      <xsl:otherwise>
        <xsl:value-of select="$text"/>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>
</xsl:stylesheet>
