using System;
using System.Collections.Generic;

namespace GeneratePath.Model
{
    public class Bridge : IEquatable<Bridge>
    {
        public List<int> Elements { get; set; }

        public override string ToString()
        {
            return "new Bridge( new List<int> {" + Utils.PathUtils.AddCommas(Elements) + "})";
        }


        public bool Equals(Bridge other)
        {
            if (ReferenceEquals(null, other)) return false;
            return ReferenceEquals(this, other) || Equals(Elements, other.Elements);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            return obj.GetType() == this.GetType() && Equals((Bridge) obj);
        }

        public override int GetHashCode() => Elements?.GetHashCode() ?? 0;

        public static bool operator ==(Bridge left, Bridge right) => Equals(left, right);

        public static bool operator !=(Bridge left, Bridge right) => !Equals(left, right);
    }
}